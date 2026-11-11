import express from "express";
import session from "express-session";
const app=express();
export function replacedCookie(req,res){return session({secret:"controlled-corpus-only",cookie:{secure:true,httpOnly:true},...req.body.options})(req,res,()=>res.end());}
export function automaticCookie(req,res){return session({secret:"controlled-corpus-only",cookie:{secure:"auto",httpOnly:true}})(req,res,()=>res.end());}
export function dynamicCookieKey(req,res){return session({secret:"controlled-corpus-only",cookie:{secure:true,[req.body.key]:req.body.value}})(req,res,()=>res.end());}
app.post("/replace",replacedCookie);app.post("/auto",automaticCookie);app.post("/computed",dynamicCookieKey);
export default app;
