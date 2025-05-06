const std = @import("std");

const Task = @This();

name: [2048]u8,
name_len: usize,
outputDir: std.fs.Dir,
outputDirPath: []const u8,

pub fn init(self: *Task, inputPath: []u8, outputDirPath: []const u8, outputDir: std.fs.Dir) !void {
    // can't use path directly (as it's a pointer owned by the caller and
    // could be deallocated before this Task is deallocated).
    if (inputPath.len > self.name.len) {
        return error.PathTooLong;
    }
    self.name_len = inputPath.len;
    std.mem.copyForwards(u8, self.name[0..inputPath.len], inputPath);

    self.outputDir = outputDir;
    self.outputDirPath = outputDirPath;
}

pub fn getName(self: *const Task) []const u8 {
    return self.name[0..self.name_len];
}

pub fn mkdir(self: *const Task, name: []const u8) !void {
    try self.outputDir.makePath(name);
}

pub fn open(self: *const Task) !std.fs.File {
    return std.fs.openFileAbsolute(self.getName(), .{});
}
