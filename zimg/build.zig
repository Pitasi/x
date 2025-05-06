const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize_mode = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "zimg",
        .root_source_file = b.path("main.zig"),
        .target = target,
        .optimize = optimize_mode,
    });

    // Link against ImageMagick libraries
    exe.linkSystemLibrary("MagickWand");
    exe.linkSystemLibrary("MagickCore");

    b.installArtifact(exe);
}
