import express from "express";
import session from "express-session";
import cors from "cors";
import Tokens from "csrf";
const app=express();const tokens=new Tokens();
app.use(session({secret:"controlled-corpus-only",resave:false,saveUninitialized:false,cookie:{secure:true,httpOnly:true}}));
export function csrfDirect(req:any,res:any){if(!tokens.verify(req.session.csrfSecret,req.headers["x-csrf-token"])) {throw new Error("invalid");} req.session.counter=req.body.counter;return res.json({ok:true});}
app.post("/csrfDirect",csrfDirect);
export function csrfNumeric(req:any,res:any){if(+tokens.verify(req.session.csrfSecret,req.headers["x-csrf-token"])) {throw new Error("invalid");} req.session.counter=req.body.counter;return res.json({ok:true});}
app.post("/csrfNumeric",csrfNumeric);
export function csrfConditional(req:any,res:any){if(req.body.validate){if(!tokens.verify(req.session.csrfSecret,req.headers["x-csrf-token"])) {throw new Error("invalid");}} req.session.counter=req.body.counter;return res.json({ok:true});}
app.post("/csrfConditional",csrfConditional);
export function csrfLate(req:any,res:any){req.session.counter=req.body.counter;if(!tokens.verify(req.session.csrfSecret,req.headers["x-csrf-token"])) {throw new Error("invalid");}return res.json({ok:true});}
app.post("/csrfLate",csrfLate);
export function csrfWrongSecret(req:any,res:any){if(!tokens.verify(req.body.secret,req.headers["x-csrf-token"])) {throw new Error("invalid");}req.session.counter=req.body.counter;return res.json({ok:true});}
app.post("/csrfWrongSecret",csrfWrongSecret);
export function redirectConditional(req:any,res:any){const target=req.query.target;if(req.body.validate){if(!target.startsWith("/account/")){throw new Error("invalid");}}return res.redirect(target);}
app.post("/redirectConditional",redirectConditional);
export function redirectWrongValue(req:any,res:any){const checked=req.query.checked;if(!checked.startsWith("/account/")){throw new Error("invalid");}return res.redirect(req.query.target);}
app.post("/redirectWrongValue",redirectWrongValue);
export function redirectReassigned(req:any,res:any){if(!req.query.target.startsWith("/account/")){throw new Error("invalid");}req.query.target=req.body.target;return res.redirect(req.query.target);}
app.post("/redirectReassigned",redirectReassigned);
export function originArrayControl(){return cors({credentials:true,origin:(origin,callback)=>callback(null,["https://trusted.example"].includes(origin))});}
app.use(originArrayControl());
export function originStringQuestion(){return cors({credentials:true,origin:(origin,callback)=>callback(null,"https://trusted.example".includes(origin))});}
app.use(originStringQuestion());
export function originNumericQuestion(){return cors({credentials:true,origin:1});}
app.use(originNumericQuestion());
export default app;
