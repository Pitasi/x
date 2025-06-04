const std = @import("std");
const Task = @import("Task.zig");
const Allocator = std.mem.Allocator;
const sha256 = std.crypto.hash.sha2.Sha256;

const c = @cImport({
    @cInclude("MagickWand/MagickWand.h");
});
