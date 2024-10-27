"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isUrlString = exports.isNonEmptyString = exports.isString = void 0;
function isString(str) {
    return typeof str === "string";
}
exports.isString = isString;
function isNonEmptyString(str) {
    return !(str === undefined || str === null || str.length === 0);
}
exports.isNonEmptyString = isNonEmptyString;
function isUrlString(str) {
    if (isNonEmptyString(str)) {
        try {
            new URL(str);
            return true;
        }
        catch (_a) {
            // No-op
        }
    }
    return false;
}
exports.isUrlString = isUrlString;
