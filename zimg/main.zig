const std = @import("std");
const Task = @import("Task.zig");
const Allocator = std.mem.Allocator;
const sha256 = std.crypto.hash.sha2.Sha256;

const c = @cImport({
    @cInclude("MagickWand/MagickWand.h");
});

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer {
        const deinit_status = gpa.deinit();
        if (deinit_status == .leak) @panic("memory leaked");
    }

    c.MagickWandGenesis();
    defer c.MagickWandTerminus();

    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len != 3) {
        std.debug.print("wrong args number, must be 3", .{});
        return;
    }

    const inputDir = args[1];
    const outputDirPath = args[2];

    try std.fs.cwd().makePath(outputDirPath);

    var outputDir = try std.fs.cwd().openDir(outputDirPath, .{});
    defer outputDir.close();

    var wg: std.Thread.WaitGroup = .{};

    const cpuCount = try std.Thread.getCpuCount();
    var pool: std.Thread.Pool = undefined;
    try pool.init(.{
        .allocator = allocator,
        .n_jobs = @truncate(cpuCount - 1),
    });
    defer pool.deinit();

    var d = try std.fs.cwd().openDir(inputDir, .{});
    defer d.close();

    var walker = try d.walk(allocator);
    defer walker.deinit();
    while (try walker.next()) |entry| {
        if (entry.kind != .file or !isPicture(entry.path)) {
            continue;
        }
        const task = try allocator.create(Task);
        const inputPath = try std.fs.path.joinZ(allocator, &[_][]const u8{ inputDir, entry.path });
        defer allocator.free(inputPath);
        try task.init(inputPath, outputDirPath, outputDir);
        pool.spawnWg(&wg, work, .{ allocator, task });
    }

    pool.waitAndWork(&wg);
}

pub fn hashfile(data: []const u8, out: *[64]u8) !void {
    var hash: [32]u8 = undefined;
    sha256.hash(data, &hash, .{});

    const slice = try std.fmt.bufPrint(out, "{}", .{
        std.fmt.fmtSliceHexLower(&hash),
    });
    if (slice.len != 64) {
        @panic("something went wrong: hexadecimal hash must be 64 bytes");
    }
}

pub fn work(allocator: Allocator, task: *Task) void {
    defer allocator.destroy(task);
    innerWork(allocator, task) catch |err| {
        std.debug.print("error processing {s}: {any}\n", .{ task.getName(), err });
    };
}

const Jobs = packed struct {
    info: bool,
    orig: bool,
    q_orig: bool,
    q_500: bool,
    q_1000: bool,
    q_2000: bool,
    w_1200: bool,
    w_1900: bool,
    w_2500: bool,
    blur: bool,
};

pub fn innerWork(allocator: Allocator, task: *Task) !void {
    const f = try task.open();
    defer f.close();

    // read file into memory
    const data = try f.readToEndAlloc(allocator, 100 * 1024 * 1024);
    defer allocator.free(data);

    // sha256 of content
    var hash: [64]u8 = undefined;
    try hashfile(data, &hash);

    // prepare output dir
    try task.mkdir(hash[0..]);

    // gather jobs to do
    const outJsonPath = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "info.json" });
    defer allocator.free(outJsonPath);

    const outOrigPath = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "orig.webp" });
    defer allocator.free(outOrigPath);

    const outSquareOrigPath = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "q_orig.webp" });
    defer allocator.free(outSquareOrigPath);

    const blurPath = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "blur.webp" });
    defer allocator.free(blurPath);

    const outQ500Path = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "q_500.webp" });
    defer allocator.free(outQ500Path);

    const outQ1000Path = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "q_1000.webp" });
    defer allocator.free(outQ1000Path);

    const outQ2000Path = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "q_2000.webp" });
    defer allocator.free(outQ2000Path);

    const outW1200Path = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "w_1200.webp" });
    defer allocator.free(outW1200Path);

    const outW1900Path = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "w_1900.webp" });
    defer allocator.free(outW1900Path);

    const outW2500Path = try std.fs.path.joinZ(allocator, &[_][]const u8{ task.outputDirPath, hash[0..], "w_2500.webp" });
    defer allocator.free(outW2500Path);

    const jobs = Jobs{
        .info = !fileExistsZ(outJsonPath),
        .blur = !fileExistsZ(blurPath),
        .orig = !fileExistsZ(outOrigPath),
        .q_orig = !fileExistsZ(outSquareOrigPath),
        .q_500 = !fileExistsZ(outQ500Path),
        .q_1000 = !fileExistsZ(outQ1000Path),
        .q_2000 = !fileExistsZ(outQ2000Path),
        .w_1200 = !fileExistsZ(outW1200Path),
        .w_1900 = !fileExistsZ(outW1900Path),
        .w_2500 = !fileExistsZ(outW2500Path),
    };

    if (!(jobs.info or
        jobs.blur or
        jobs.orig or
        jobs.q_orig or
        jobs.q_500 or
        jobs.q_1000 or
        jobs.q_2000 or
        jobs.w_1200 or
        jobs.w_1900 or
        jobs.w_2500))
    {
        return;
    }

    // init imagemagick wand
    const wand = c.NewMagickWand();
    if (wand == null) {
        std.debug.print("Failed to create MagickWand\n", .{});
        return;
    }
    defer _ = c.DestroyMagickWand(wand);

    if (c.MagickReadImageBlob(wand, data.ptr, data.len) == c.MagickFalse) {
        std.debug.print("Failed to read input image\n", .{});
        return;
    }

    const width = c.MagickGetImageWidth(wand);
    const height = c.MagickGetImageHeight(wand);

    // init imagemagick wand and crop to square
    const squareWand = c.CloneMagickWand(wand);
    defer _ = c.DestroyMagickWand(squareWand);
    if (width != height) {
        const size = @min(width, height);
        const x: isize = if (width > height) @intCast((width - size) / 2) else 0;
        const y: isize = if (width < height) @intCast((height - size) / 2) else 0;

        if (c.MagickCropImage(squareWand, size, size, x, y) == c.MagickFalse) {
            std.debug.print("Failed to crop image\n", .{});
            return;
        }
    }

    if (jobs.info) {
        const tmpWand = c.CloneMagickWand(wand);
        defer _ = c.DestroyMagickWand(tmpWand);
        if (c.MagickWriteImage(tmpWand, outJsonPath.ptr) == c.MagickFalse) {
            std.debug.print("Failed to write output image to {s}\n", .{outJsonPath});
            return;
        }

        std.debug.print("written {s}\n", .{outJsonPath});
    }

    if (c.MagickStripImage(wand) == c.MagickFalse) {
        std.debug.print("Failed to strip image", .{});
        return;
    }

    if (c.MagickStripImage(squareWand) == c.MagickFalse) {
        std.debug.print("Failed to strip image", .{});
        return;
    }

    if (jobs.orig) {
        const tmpWand = c.CloneMagickWand(wand);
        defer _ = c.DestroyMagickWand(tmpWand);
        if (c.MagickWriteImage(tmpWand, outOrigPath.ptr) == c.MagickFalse) {
            std.debug.print("Failed to write output image to {s}\n", .{outOrigPath});
            return;
        }
        std.debug.print("written {s}\n", .{outOrigPath});
    }

    if (jobs.q_orig) {
        const tmpWand = c.CloneMagickWand(squareWand);
        defer _ = c.DestroyMagickWand(tmpWand);

        if (c.MagickWriteImage(tmpWand, outSquareOrigPath.ptr) == c.MagickFalse) {
            std.debug.print("Failed to write output image to {s}\n", .{outSquareOrigPath});
            return;
        }
        std.debug.print("written {s}\n", .{outSquareOrigPath});
    }

    if (jobs.blur) {
        const tmpWand = c.CloneMagickWand(wand);
        defer _ = c.DestroyMagickWand(tmpWand);

        const maxWidth: f64 = 25;

        const factor: f64 = maxWidth / @as(f64, @floatFromInt(width));
        const targetHeight: f64 = @as(f64, @floatFromInt(height)) * factor;

        const targetWidth: u64 = @intFromFloat(@floor(maxWidth));
        const targetHeightTrunc: u64 = @intFromFloat(@floor(targetHeight));

        if (c.MagickResizeImage(tmpWand, targetWidth, targetHeightTrunc, c.LanczosFilter) == c.MagickFalse) {
            std.debug.print("Failed to crop image\n", .{});
            return;
        }

        if (c.MagickSetImageCompressionQuality(tmpWand, 5) == c.MagickFalse) {
            std.debug.print("Failed to set image compression quality\n", .{});
            return;
        }

        if (c.MagickWriteImage(tmpWand, blurPath.ptr) == c.MagickFalse) {
            std.debug.print("Failed to write output image to {s}\n", .{blurPath});
            return;
        }
        std.debug.print("written {s}\n", .{blurPath});
    }

    if (jobs.q_500) {
        try q(squareWand, width, height, outQ500Path, 500);
    }

    if (jobs.q_1000 and @min(width, height) > 1000) {
        try q(squareWand, width, height, outQ1000Path, 1000);
    }

    if (jobs.q_2000 and @min(width, height) > 2000) {
        try q(squareWand, width, height, outQ2000Path, 2000);
    }

    if (jobs.w_1200 and width > 1200) {
        try w(wand, width, height, outW1200Path, 1200);
    }
    if (jobs.w_1900 and width > 1900) {
        try w(wand, width, height, outW1900Path, 1900);
    }
    if (jobs.w_2500 and width > 2500) {
        try w(wand, width, height, outW2500Path, 2500);
    }
}

pub fn isPicture(path: []const u8) bool {
    const e = std.fs.path.extension(path);
    return std.mem.eql(u8, e, ".jpg") or
        std.mem.eql(u8, e, ".jpeg") or
        std.mem.eql(u8, e, ".png");
}

fn q(squareWand: ?*c.MagickWand, width: u64, height: u64, targetFile: [:0]const u8, targetSize: u64) !void {
    const tmpWand = c.CloneMagickWand(squareWand);
    defer _ = c.DestroyMagickWand(tmpWand);

    const size = @min(width, height);
    const qSize = @min(size, targetSize);
    if (c.MagickResizeImage(tmpWand, qSize, qSize, c.LanczosFilter) == c.MagickFalse) {
        std.debug.print("Failed to crop image\n", .{});
        return;
    }

    if (c.MagickSetImageCompressionQuality(tmpWand, 10) == c.MagickFalse) {
        std.debug.print("Failed to set image compression quality\n", .{});
        return;
    }

    if (c.MagickWriteImage(tmpWand, targetFile.ptr) == c.MagickFalse) {
        std.debug.print("Failed to write output image to {s}\n", .{targetFile});
        return;
    }
    std.debug.print("written {s}\n", .{targetFile});
}

fn w(wand: ?*c.MagickWand, width: u64, height: u64, targetFile: []const u8, maxWidth: f64) !void {
    const tmpWand = c.CloneMagickWand(wand);
    defer _ = c.DestroyMagickWand(tmpWand);

    const factor: f64 = maxWidth / @as(f64, @floatFromInt(width));
    const targetHeight: f64 = @as(f64, @floatFromInt(height)) * factor;

    const targetWidth: u64 = @intFromFloat(@floor(maxWidth));
    const targetHeightTrunc: u64 = @intFromFloat(@floor(targetHeight));

    if (c.MagickResizeImage(tmpWand, targetWidth, targetHeightTrunc, c.LanczosFilter) == c.MagickFalse) {
        std.debug.print("Failed to crop image\n", .{});
        return;
    }

    if (c.MagickSetImageCompressionQuality(tmpWand, 80) == c.MagickFalse) {
        std.debug.print("Failed to set image compression quality\n", .{});
        return;
    }

    if (c.MagickWriteImage(tmpWand, targetFile.ptr) == c.MagickFalse) {
        std.debug.print("Failed to write output image to {s}\n", .{targetFile});
        return;
    }
    std.debug.print("written {s}\n", .{targetFile});
}

fn fileExistsZ(path: [:0]u8) bool {
    const access = std.fs.accessAbsoluteZ(path, .{});
    if (access) |_| {
        return true;
    } else |_| {
        return false;
    }
}
