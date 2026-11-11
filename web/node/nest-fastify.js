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
sessionControl(req,res) { return res.setCookie("sid",req.session.id,{secure:true,httpOnly:true}).send({ok:true}); }
sessionPositive(req,res) { return res.setCookie("sid",req.session.id,{secure:false,httpOnly:false}).send({ok:true}); }
  htmlPositive(req, res) { return res.type("html").send(req.query.html); }
  htmlControl(req, res) { return res.send({ value: req.query.html }); }
  htmlUncertain(req, res) { return res.type("html").send(fragment(req.query.html)); }
  csrfPositive(req, res) { req.session.counter = req.body.counter; return res.send({ ok: true }); }
  csrfControl(req, res) { if (!tokens.verify(req.session.csrfSecret, req.headers["x-csrf-token"])) { throw new Error("invalid token"); } req.session.counter = req.body.counter; return res.send({ ok: true }); }
  csrfUncertain(req, res) { req.session.counter = req.body.counter; return res.send({ ok: true }); }
  async authPositive(req, res) { if (req.user.role !== "writer") { res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  async authControl(req, res) { if (req.user.role !== "writer") { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  async authUncertain(req, res) { if (req.user.role !== "writer" && !enabled(req)) { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  redirectPositive(req, res) { return res.redirect(req.query.target); }
  redirectControl(req, res) { const target = req.query.target; if (!target.startsWith("/account/")) { return res.status(400).send("bad target"); } return res.redirect(target); }
  redirectUncertain(req, res) { return res.redirect(fragment(req.query.target)); }
  async massPositive(req, res) { return prisma.user.create({ data: { tenant: req.user.tenant, ...req.body } }); }
  async massControl(req, res) { return prisma.user.create({ data: { name: req.body.name, tenant: req.user.tenant } }); }
  async massUncertain(req, res) { return prisma.user.create({ data: { tenant: req.user.tenant, ...transform(req.body) } }); }
  corsPositive(req, res) { res.header("Access-Control-Allow-Credentials", "true"); res.header("Access-Control-Allow-Origin", req.headers.origin); return res.send({ ok: true }); }
  corsControl(req, res) { res.header("Access-Control-Allow-Credentials", "true"); res.header("Access-Control-Allow-Origin", "https://trusted.example"); return res.send({ ok: true }); }
  corsUncertain(req, res) { res.header("Access-Control-Allow-Credentials", "true"); res.header("Access-Control-Allow-Origin", policyOrigin(req.headers.origin)); return res.send({ ok: true }); }
  sessionUncertain(req, res) { return res.setCookie("sid", req.session.id, {secure:enabled(req),httpOnly:true}).send({ok:true}); }
  tokenIssue(req,res) { req.session.csrfSecret=tokens.secretSync(); return res.send({token:tokens.create(req.session.csrfSecret)}); }
}
Post("htmlPositive")(CorpusController.prototype,"htmlPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"htmlPositive"));
Req()(CorpusController.prototype,"htmlPositive",0);
Res()(CorpusController.prototype,"htmlPositive",1);
Post("htmlControl")(CorpusController.prototype,"htmlControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"htmlControl"));
Req()(CorpusController.prototype,"htmlControl",0);
Res()(CorpusController.prototype,"htmlControl",1);
Post("htmlUncertain")(CorpusController.prototype,"htmlUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"htmlUncertain"));
Req()(CorpusController.prototype,"htmlUncertain",0);
Res()(CorpusController.prototype,"htmlUncertain",1);
Post("csrfPositive")(CorpusController.prototype,"csrfPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"csrfPositive"));
Req()(CorpusController.prototype,"csrfPositive",0);
Res()(CorpusController.prototype,"csrfPositive",1);
Post("csrfControl")(CorpusController.prototype,"csrfControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"csrfControl"));
Req()(CorpusController.prototype,"csrfControl",0);
Res()(CorpusController.prototype,"csrfControl",1);
Post("csrfUncertain")(CorpusController.prototype,"csrfUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"csrfUncertain"));
Req()(CorpusController.prototype,"csrfUncertain",0);
Res()(CorpusController.prototype,"csrfUncertain",1);
Post("authPositive")(CorpusController.prototype,"authPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"authPositive"));
Req()(CorpusController.prototype,"authPositive",0);
Res()(CorpusController.prototype,"authPositive",1);
Post("authControl")(CorpusController.prototype,"authControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"authControl"));
Req()(CorpusController.prototype,"authControl",0);
Res()(CorpusController.prototype,"authControl",1);
Post("authUncertain")(CorpusController.prototype,"authUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"authUncertain"));
Req()(CorpusController.prototype,"authUncertain",0);
Res()(CorpusController.prototype,"authUncertain",1);
Post("redirectPositive")(CorpusController.prototype,"redirectPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"redirectPositive"));
Req()(CorpusController.prototype,"redirectPositive",0);
Res()(CorpusController.prototype,"redirectPositive",1);
Post("redirectControl")(CorpusController.prototype,"redirectControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"redirectControl"));
Req()(CorpusController.prototype,"redirectControl",0);
Res()(CorpusController.prototype,"redirectControl",1);
Post("redirectUncertain")(CorpusController.prototype,"redirectUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"redirectUncertain"));
Req()(CorpusController.prototype,"redirectUncertain",0);
Res()(CorpusController.prototype,"redirectUncertain",1);
Post("massPositive")(CorpusController.prototype,"massPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"massPositive"));
Req()(CorpusController.prototype,"massPositive",0);
Res()(CorpusController.prototype,"massPositive",1);
Post("massControl")(CorpusController.prototype,"massControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"massControl"));
Req()(CorpusController.prototype,"massControl",0);
Res()(CorpusController.prototype,"massControl",1);
Post("massUncertain")(CorpusController.prototype,"massUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"massUncertain"));
Req()(CorpusController.prototype,"massUncertain",0);
Res()(CorpusController.prototype,"massUncertain",1);
Post("corsPositive")(CorpusController.prototype,"corsPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"corsPositive"));
Req()(CorpusController.prototype,"corsPositive",0);
Res()(CorpusController.prototype,"corsPositive",1);
Post("corsControl")(CorpusController.prototype,"corsControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"corsControl"));
Req()(CorpusController.prototype,"corsControl",0);
Res()(CorpusController.prototype,"corsControl",1);
Post("corsUncertain")(CorpusController.prototype,"corsUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"corsUncertain"));
Req()(CorpusController.prototype,"corsUncertain",0);
Res()(CorpusController.prototype,"corsUncertain",1);
Post("sessionUncertain")(CorpusController.prototype,"sessionUncertain",Object.getOwnPropertyDescriptor(CorpusController.prototype,"sessionUncertain"));
Req()(CorpusController.prototype,"sessionUncertain",0);
Res()(CorpusController.prototype,"sessionUncertain",1);
Post("tokenIssue")(CorpusController.prototype,"tokenIssue",Object.getOwnPropertyDescriptor(CorpusController.prototype,"tokenIssue"));
Req()(CorpusController.prototype,"tokenIssue",0);
Res()(CorpusController.prototype,"tokenIssue",1);
Controller("corpus")(CorpusController);
export class CorpusModule {}
Module({controllers:[CorpusController]})(CorpusModule);
export async function bootstrapPositive() {
const app=await NestFactory.create(CorpusModule,new FastifyAdapter());
await app.register(fastifyCookie);
await app.register(fastifySession,{secret:"controlled-corpus-secret-32-characters",cookie:{secure:false,httpOnly:true}});
return app;
}
export async function bootstrapControl() {
const app=await NestFactory.create(CorpusModule,new FastifyAdapter());
await app.register(fastifyCookie);
await app.register(fastifySession,{secret:"controlled-corpus-secret-32-characters",cookie:{secure:true,httpOnly:true}});
return app;
}
export async function bootstrapUncertain() {
const app=await NestFactory.create(CorpusModule,new FastifyAdapter());
await app.register(fastifyCookie);
await app.register(fastifySession,{secret:"controlled-corpus-secret-32-characters",cookie:{secure:enabled(process.env),httpOnly:true}});
return app;
}
Post("sessionPositive")(CorpusController.prototype,"sessionPositive",Object.getOwnPropertyDescriptor(CorpusController.prototype,"sessionPositive"));
Req()(CorpusController.prototype,"sessionPositive",0);
Res()(CorpusController.prototype,"sessionPositive",1);
Post("sessionControl")(CorpusController.prototype,"sessionControl",Object.getOwnPropertyDescriptor(CorpusController.prototype,"sessionControl"));
Req()(CorpusController.prototype,"sessionControl",0);
Res()(CorpusController.prototype,"sessionControl",1);
export class MetadataController {
publicOperation(){return "observed";}
requiredOperation(){return "observed";}
}
ReadAccess(Modes.open)(MetadataController.prototype,"publicOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"publicOperation"));
Post("publicOperation")(MetadataController.prototype,"publicOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"publicOperation"));
ReadAccess(Modes.writer)(MetadataController.prototype,"requiredOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"requiredOperation"));
Post("requiredOperation")(MetadataController.prototype,"requiredOperation",Object.getOwnPropertyDescriptor(MetadataController.prototype,"requiredOperation"));
Controller("metadata")(MetadataController);
export class MetadataModule{}
Module({controllers:[MetadataController],providers:[{provide:APP_GUARD,useClass:SourceGuard}]})(MetadataModule);
