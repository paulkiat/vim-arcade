"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const ava_1 = __importDefault(require("ava"));
const strings_js_1 = require("./strings.js");
(0, ava_1.default)("isNonEmptyString", (t) => {
    const fails = [undefined, null, ""];
    const passes = ["1", "234"];
    for (let x of fails) {
        t.is((0, strings_js_1.isNonEmptyString)(x), false);
    }
    for (let x of passes) {
        t.is((0, strings_js_1.isNonEmptyString)(x), true);
    }
});
(0, ava_1.default)("isUrlString", (t) => {
    const fails = [undefined, null, "", "https", "moov.io"];
    const passes = ["https://moov.io", "https://docs.moov.io/api/#hash"];
    for (let x of fails) {
        t.is((0, strings_js_1.isUrlString)(x), false);
    }
    for (let x of passes) {
        t.is((0, strings_js_1.isUrlString)(x), true);
    }
});
