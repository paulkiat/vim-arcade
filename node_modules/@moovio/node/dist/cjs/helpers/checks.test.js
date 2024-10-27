"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const ava_1 = __importDefault(require("ava"));
const checks_js_1 = require("./checks.js");
(0, ava_1.default)("check", (t) => {
    const throws = [undefined, null];
    const notThrows = [false, true, "", {}];
    for (let x of throws) {
        t.throws(() => (0, checks_js_1.check)(x).or());
    }
    for (let x of notThrows) {
        t.notThrows(() => (0, checks_js_1.check)(x).or());
    }
});
(0, ava_1.default)("checkString", (t) => {
    const throws = [undefined, null, ""];
    const notThrows = ["a", "1"];
    for (let x of throws) {
        t.throws(() => (0, checks_js_1.checkString)(x).or());
    }
    for (let x of notThrows) {
        t.notThrows(() => (0, checks_js_1.checkString)(x).or());
    }
});
(0, ava_1.default)("checkUrl", (t) => {
    const throws = [undefined, null, "", "http", "moov.io"];
    const notThrows = ["https://moov.io", "https://docs.moov.io/api/#hash"];
    for (let x of throws) {
        t.throws(() => (0, checks_js_1.checkUrl)(x).or());
    }
    for (let x of notThrows) {
        t.notThrows(() => (0, checks_js_1.checkUrl)(x).or());
    }
});
(0, ava_1.default)("checkArray", (t) => {
    const throws = [undefined, null, {}];
    const notThrows = [[], [1], ["one", 2]];
    for (let x of throws) {
        t.throws(() => (0, checks_js_1.checkArray)(x).or());
    }
    for (let x of notThrows) {
        t.notThrows(() => (0, checks_js_1.checkArray)(x).or());
    }
});
(0, ava_1.default)("checkArrayLength", (t) => {
    const throws = [undefined, null, {}, []];
    const notThrows = [[1], ["one", 2]];
    for (let x of throws) {
        t.throws(() => (0, checks_js_1.checkArrayLength)(x).or());
    }
    for (let x of notThrows) {
        t.notThrows(() => (0, checks_js_1.checkArrayLength)(x).or());
    }
});
