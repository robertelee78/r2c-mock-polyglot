import express from "express";
import session from "express-session";
import cors from "cors";
import {parse} from "yaml";
import {readFileSync} from "node:fs";
import programSettings from "./settings-program";
const app=express();
const yamlSettings=parse(readFileSync("web/node/settings.yaml","utf8"));
const jsonSettings=JSON.parse(readFileSync("web/node/settings.json","utf8"));
export function yamlConsumer(req:any,res:any){return session(yamlSettings)(req,res,()=>res.end());}
export function jsonConsumer(req:any,res:any){return session(jsonSettings)(req,res,()=>res.end());}
export function programConsumer(req:any,res:any){return session(programSettings.weakening)(req,res,()=>res.end());}
export function programControl(req:any,res:any){return session(programSettings.hardened)(req,res,()=>res.end());}
export function selectedProgram(req:any,res:any){return session(programSettings[req.query.profile])(req,res,()=>res.end());}
app.use(cors(yamlSettings.cors));
app.post("/yaml",yamlConsumer);app.post("/json",jsonConsumer);app.post("/program",programConsumer);app.post("/program-control",programControl);app.post("/selected",selectedProgram);
export default app;
