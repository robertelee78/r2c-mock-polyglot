export function fragment(value) { return String(value); }
export function transform(value) { return { ...value }; }
export function enabled(req) { return req.headers["x-policy"] === "enabled"; }
export function policyOrigin(origin) { return origin; }
