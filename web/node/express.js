import express from "express";
import session from "express-session";
import cors from "cors";
import csrf from "csurf";
import Handlebars from "handlebars";
import { readFile, readFileSync } from "node:fs";
import settings from "./settings.json" assert {type:"json"};
import { PrismaClient } from "@prisma/client";
import Tokens from "csrf";
import { policyOrigin, fragment, transform, enabled } from "./helpers.js";
const prisma = new PrismaClient();
const tokens = new Tokens();
const app=express();
app.use(express.json());
app.use(session({secret:"controlled-corpus-only",resave:false,saveUninitialized:false,cookie:{secure:true,httpOnly:true}}));
app.set("views","web/node/views");
app.set("view engine","ejs");
app.engine("hbs",function handlebarsEngine(file,options,done){readFile(file,"utf8",function readTemplate(error,contents){if(error)return done(error);try{done(null,Handlebars.compile(contents)(options));}catch(error){done(error);}});});
Handlebars.registerPartial("fragment",Handlebars.compile(readFileSync("web/node/views/raw.hbs","utf8")));
export function htmlPositive(req, res) { return res.type("html").send(req.query.html); }
export function htmlControl(req, res) { return res.json({ value: req.query.html }); }
export function htmlUncertain(req, res) { return res.type("html").send(fragment(req.query.html)); }
export function csrfPositive(req, res) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
export function csrfControl(req, res) { if (!tokens.verify(req.session.csrfSecret, req.headers["x-csrf-token"])) { throw new Error("invalid token"); } req.session.counter = req.body.counter; return res.send({ ok: true }); }
export function csrfUncertain(req, res) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
export async function authPositive(req, res) { if (req.user.role !== "writer") { res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
export async function authControl(req, res) { if (req.user.role !== "writer") { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
export async function authUncertain(req, res) { if (req.user.role !== "writer" && !enabled(req)) { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
export function redirectPositive(req, res) { return res.redirect(req.query.target); }
export function redirectControl(req, res) { const target = req.query.target; if (!target.startsWith("/account/")) { return res.status(400).send("bad target"); } return res.redirect(target); }
export function redirectUncertain(req, res) { return res.redirect(fragment(req.query.target)); }
export async function massPositive(req, res) { return prisma.user.create({ data: { tenant: req.user.tenant, ...req.body } }); }
export async function massControl(req, res) { return prisma.user.create({ data: { name: req.body.name, tenant: req.user.tenant } }); }
export async function massUncertain(req, res) { return prisma.user.create({ data: { tenant: req.user.tenant, ...transform(req.body) } }); }
export function corsPositive(req, res) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", req.headers.origin); return res.json({ ok: true }); }
export function corsControl(req, res) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", "https://trusted.example"); return res.json({ ok: true }); }
export function corsUncertain(req, res) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", policyOrigin(req.headers.origin)); return res.json({ ok: true }); }
export function templateUncertain(req, res) { return res.render("raw.ejs", { value: fragment(req.query.value) }); }
export function sessionUncertain(req, res) { return res.cookie("sid", req.session.id, {secure:enabled(req),httpOnly:true}).send({ok:true}); }
export function tokenIssue(req,res) { req.session.csrfSecret=tokens.secretSync(); return res.json({token:tokens.create(req.session.csrfSecret)}); }
export function templateEjsPositive(req,res) { return res.render("raw.ejs", {value:req.query.value}); }
export function templateEjsControl(req,res) { return res.render("escaped.ejs", {value:req.query.value}); }
export function templateEjsUncertain(req,res) { return res.render("raw.ejs", {value:fragment(req.query.value)}); }
export function templatePugPositive(req,res) { return res.render("raw.pug", {value:req.query.value}); }
export function templatePugControl(req,res) { return res.render("escaped.pug", {value:req.query.value}); }
export function templatePugUncertain(req,res) { return res.render("raw.pug", {value:fragment(req.query.value)}); }
export function templateHbsPositive(req,res) { return res.render("raw.hbs", {value:req.query.value}); }
export function templateHbsControl(req,res) { return res.render("escaped.hbs", {value:req.query.value}); }
export function templateHbsUncertain(req,res) { return res.render("raw.hbs", {value:fragment(req.query.value)}); }
app.post("/htmlPositive", htmlPositive);
app.post("/htmlControl", htmlControl);
app.post("/htmlUncertain", htmlUncertain);
app.post("/csrfPositive", csrfPositive);
app.post("/csrfControl", csrf(), csrfControl);
app.post("/csrfUncertain", csrfUncertain);
app.post("/authPositive", authPositive);
app.post("/authControl", authControl);
app.post("/authUncertain", authUncertain);
app.post("/redirectPositive", redirectPositive);
app.post("/redirectControl", redirectControl);
app.post("/redirectUncertain", redirectUncertain);
app.post("/massPositive", massPositive);
app.post("/massControl", massControl);
app.post("/massUncertain", massUncertain);
app.post("/corsPositive", corsPositive);
app.post("/corsControl", corsControl);
app.post("/corsUncertain", corsUncertain);
app.post("/templateUncertain", templateUncertain);
app.post("/sessionUncertain", sessionUncertain);
app.post("/tokenIssue", tokenIssue);
app.post("/templateEjsPositive", templateEjsPositive);
app.post("/templateEjsControl", templateEjsControl);
app.post("/templateEjsUncertain", templateEjsUncertain);
app.post("/templatePugPositive", templatePugPositive);
app.post("/templatePugControl", templatePugControl);
app.post("/templatePugUncertain", templatePugUncertain);
app.post("/templateHbsPositive", templateHbsPositive);
app.post("/templateHbsControl", templateHbsControl);
app.post("/templateHbsUncertain", templateHbsUncertain);
export function settingsConsumer(req,res){session(settings)(req,res,()=>res.end());}
app.post("/settings",settingsConsumer);
export function sessionPositive(req,res) { return res.cookie("sid",req.session.id,{secure:false,httpOnly:false}).send({ok:true}); }
app.post("/sessionPositive",sessionPositive);
export function sessionControl(req,res) { return res.cookie("sid",req.session.id,{secure:true,httpOnly:true}).send({ok:true}); }
app.post("/sessionControl",sessionControl);
export default app;
