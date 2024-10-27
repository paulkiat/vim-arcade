"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.checkArrayLength = exports.checkArray = exports.checkUrl = exports.checkString = exports.check = void 0;
const strings_js_1 = require("./strings.js");
const orThrow = {
    or: (msg) => {
        throw new TypeError(msg);
    },
};
const orNothing = {
    or: () => { },
};
function check(any) {
    return !(any === undefined || any === null) ? orNothing : orThrow;
}
exports.check = check;
function checkString(str) {
    return (0, strings_js_1.isNonEmptyString)(str) ? orNothing : orThrow;
}
exports.checkString = checkString;
function checkUrl(url) {
    return (0, strings_js_1.isUrlString)(url) ? orNothing : orThrow;
}
exports.checkUrl = checkUrl;
function checkArray(arr) {
    return check(arr) && arr.constructor === Array ? orNothing : orThrow;
}
exports.checkArray = checkArray;
function checkArrayLength(arr) {
    return checkArray(arr) && arr.length > 0 ? orNothing : orThrow;
}
exports.checkArrayLength = checkArrayLength;
