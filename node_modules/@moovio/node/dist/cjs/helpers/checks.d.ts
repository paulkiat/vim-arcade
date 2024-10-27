export function check(any: any): {
    or: (msg: any) => never;
} | {
    or: () => void;
};
export function checkString(str: any): {
    or: (msg: any) => never;
} | {
    or: () => void;
};
export function checkUrl(url: any): {
    or: (msg: any) => never;
} | {
    or: () => void;
};
export function checkArray(arr: any): {
    or: (msg: any) => never;
} | {
    or: () => void;
};
export function checkArrayLength(arr: any): {
    or: (msg: any) => never;
} | {
    or: () => void;
};
