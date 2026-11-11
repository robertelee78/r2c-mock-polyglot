import {ReadAccess} from "./access.js";
import {SourceGuard} from "./source-guard.js";
import {Modes} from "./permissions.js";
import { Controller,Post,Req,Res,Body,Module,SetMetadata } from "@nestjs/common";
import { NestFactory,Reflector,APP_GUARD } from "@nestjs/core";
import { Mutation,Resolver,Context } from "@nestjs/graphql";
import { FastifyAdapter } from "@nestjs/platform-fastify";
import session from "express-session";
import fastifySession from "@fastify/session";
import fastifyCookie from "@fastify/cookie";
import { PrismaClient } from "@prisma/client";
import Tokens from "csrf";
import { policyOrigin, fragment, transform, enabled } from "./helpers.js";
const prisma = new PrismaClient();
const tokens = new Tokens();
export class CorpusController {
sessionControl(req,res) { return res.cookie("sid",req.session.id,{secure:true,httpOnly:true}).send({ok:true}); }
sessionPositive(req,res) { return res.cookie("sid",req.session.id,{secure:false,httpOnly:false}).send({ok:true}); }
  htmlPositive(req, res) { return res.type("html").send(req.query.html); }
  htmlControl(req, res) { return res.json({ value: req.query.html }); }
  htmlUncertain(req, res) { return res.type("html").send(fragment(req.query.html)); }
  csrfPositive(req, res) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
  csrfControl(req, res) { if (!tokens.verify(req.session.csrfSecret, req.headers["x-csrf-token"])) { throw new Error("invalid token"); } req.session.counter = req.body.counter; return res.send({ ok: true }); }
  csrfUncertain(req, res) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
  async authPositive(req, res) { if (req.user.role !== "writer") { res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  async authControl(req, res) { if (req.user.role !== "writer") { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  async authUncertain(req, res) { if (req.user.role !== "writer" && !enabled(req)) { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  redirectPositive(req, res) { return res.redirect(req.query.target); }
  redirectControl(req, res) { const target = req.query.target; if (!target.startsWith("/account/")) { return res.status(400).send("bad target"); } return res.redirect(target); }
  redirectUncertain(req, res) { return res.redirect(fragment(req.query.target)); }
  async massPositive(req, res) { return prisma.user.create({ data: { tenant: req.user.tenant, ...req.body } }); }
  async massControl(req, res) { return prisma.user.create({ data: { name: req.body.name, tenant: req.user.tenant } }); }
  async massUncertain(req, res) { return prisma.user.create({ data: { tenant: req.user.tenant, ...transform(req.body) } }); }
  corsPositive(req, res) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", req.headers.origin); return res.json({ ok: true }); }
  corsControl(req, res) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", "https://trusted.example"); return res.json({ ok: true }); }
  corsUncertain(req, res) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", policyOrigin(req.headers.origin)); return res.json({ ok: true }); }
  sessionUncertain(req, res) { return res.cookie("sid", req.session.id, {secure:enabled(req),httpOnly:true}).send({ok:true}); }
  tokenIssue(req,res) { req.session.csrfSecret=tokens.secretSync(); return res.json({token:tokens.create(req.session.csrfSecret)}); }
}
Mutation(() => String)(CorpusController.prototype,"htmlPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"htmlPositive"));
Context("req")(CorpusController.prototype,"htmlPositive",0);
Context("res")(CorpusController.prototype,"htmlPositive",1);
Mutation(() => String)(CorpusController.prototype,"htmlControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"htmlControl"));
Context("req")(CorpusController.prototype,"htmlControl",0);
Context("res")(CorpusController.prototype,"htmlControl",1);
Mutation(() => String)(CorpusController.prototype,"htmlUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"htmlUncertain"));
Context("req")(CorpusController.prototype,"htmlUncertain",0);
Context("res")(CorpusController.prototype,"htmlUncertain",1);
Mutation(() => String)(CorpusController.prototype,"csrfPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"csrfPositive"));
Context("req")(CorpusController.prototype,"csrfPositive",0);
Context("res")(CorpusController.prototype,"csrfPositive",1);
Mutation(() => String)(CorpusController.prototype,"csrfControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"csrfControl"));
Context("req")(CorpusController.prototype,"csrfControl",0);
Context("res")(CorpusController.prototype,"csrfControl",1);
Mutation(() => String)(CorpusController.prototype,"csrfUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"csrfUncertain"));
Context("req")(CorpusController.prototype,"csrfUncertain",0);
Context("res")(CorpusController.prototype,"csrfUncertain",1);
Mutation(() => String)(CorpusController.prototype,"authPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"authPositive"));
Context("req")(CorpusController.prototype,"authPositive",0);
Context("res")(CorpusController.prototype,"authPositive",1);
Mutation(() => String)(CorpusController.prototype,"authControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"authControl"));
Context("req")(CorpusController.prototype,"authControl",0);
Context("res")(CorpusController.prototype,"authControl",1);
Mutation(() => String)(CorpusController.prototype,"authUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"authUncertain"));
Context("req")(CorpusController.prototype,"authUncertain",0);
Context("res")(CorpusController.prototype,"authUncertain",1);
Mutation(() => String)(CorpusController.prototype,"redirectPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"redirectPositive"));
Context("req")(CorpusController.prototype,"redirectPositive",0);
Context("res")(CorpusController.prototype,"redirectPositive",1);
Mutation(() => String)(CorpusController.prototype,"redirectControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"redirectControl"));
Context("req")(CorpusController.prototype,"redirectControl",0);
Context("res")(CorpusController.prototype,"redirectControl",1);
Mutation(() => String)(CorpusController.prototype,"redirectUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"redirectUncertain"));
Context("req")(CorpusController.prototype,"redirectUncertain",0);
Context("res")(CorpusController.prototype,"redirectUncertain",1);
Mutation(() => String)(CorpusController.prototype,"massPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"massPositive"));
Context("req")(CorpusController.prototype,"massPositive",0);
Context("res")(CorpusController.prototype,"massPositive",1);
Mutation(() => String)(CorpusController.prototype,"massControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"massControl"));
Context("req")(CorpusController.prototype,"massControl",0);
Context("res")(CorpusController.prototype,"massControl",1);
Mutation(() => String)(CorpusController.prototype,"massUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"massUncertain"));
Context("req")(CorpusController.prototype,"massUncertain",0);
Context("res")(CorpusController.prototype,"massUncertain",1);
Mutation(() => String)(CorpusController.prototype,"corsPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"corsPositive"));
Context("req")(CorpusController.prototype,"corsPositive",0);
Context("res")(CorpusController.prototype,"corsPositive",1);
Mutation(() => String)(CorpusController.prototype,"corsControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"corsControl"));
Context("req")(CorpusController.prototype,"corsControl",0);
Context("res")(CorpusController.prototype,"corsControl",1);
Mutation(() => String)(CorpusController.prototype,"corsUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"corsUncertain"));
Context("req")(CorpusController.prototype,"corsUncertain",0);
Context("res")(CorpusController.prototype,"corsUncertain",1);
Mutation(() => String)(CorpusController.prototype,"sessionUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"sessionUncertain"));
Context("req")(CorpusController.prototype,"sessionUncertain",0);
Context("res")(CorpusController.prototype,"sessionUncertain",1);
Mutation(() => String)(CorpusController.prototype,"tokenIssue",Object.getOwnPropertyDescriptor(CorpusController.prototype,"tokenIssue"));
Context("req")(CorpusController.prototype,"tokenIssue",0);
Context("res")(CorpusController.prototype,"tokenIssue",1);
Resolver()(CorpusController);
export class CorpusModule {}
Module({providers:[CorpusController]})(CorpusModule);
export async function bootstrapPositive() {
const app=await NestFactory.create(CorpusModule);
app.use(session({secret:"controlled-corpus-only",resave:false,saveUninitialized:false,cookie:{secure:false,httpOnly:true}}));
return app;
}
export async function bootstrapControl() {
const app=await NestFactory.create(CorpusModule);
app.use(session({secret:"controlled-corpus-only",resave:false,saveUninitialized:false,cookie:{secure:true,httpOnly:true}}));
return app;
}
export async function bootstrapUncertain() {
const app=await NestFactory.create(CorpusModule);
app.use(session({secret:"controlled-corpus-only",resave:false,saveUninitialized:false,cookie:{secure:enabled(process.env),httpOnly:true}}));
return app;
}
Mutation(() => String)(CorpusController.prototype,"sessionPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"sessionPositive"));
Context("req")(CorpusController.prototype,"sessionPositive",0);
Context("res")(CorpusController.prototype,"sessionPositive",1);
Mutation(() => String)(CorpusController.prototype,"sessionControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"sessionControl"));
Context("req")(CorpusController.prototype,"sessionControl",0);
Context("res")(CorpusController.prototype,"sessionControl",1);
export class MetadataController {
publicOperation(){return "observed";}
requiredOperation(){return "observed";}
}
ReadAccess(Modes.open)(MetadataController.prototype,"publicOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"publicOperation"));
Mutation(() => String)(MetadataController.prototype,"publicOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"publicOperation"));
ReadAccess(Modes.writer)(MetadataController.prototype,"requiredOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"requiredOperation"));
Mutation(() => String)(MetadataController.prototype,"requiredOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"requiredOperation"));
Resolver()(MetadataController);
export class MetadataModule{}
Module({providers:[MetadataController,{provide:APP_GUARD,useClass:SourceGuard}]})(MetadataModule);
