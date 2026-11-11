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
@Mutation(() => String)
sessionControl(@Context("req") req:any,@Context("res") res:any) { return res.cookie("sid",req.session.id,{secure:true,httpOnly:true}).send({ok:true}); }
@Mutation(() => String)
sessionPositive(@Context("req") req:any,@Context("res") res:any) { return res.cookie("sid",req.session.id,{secure:false,httpOnly:false}).send({ok:true}); }
  @Mutation(() => String)
  htmlPositive(@Context("req") req: any, @Context("res") res: any) { return res.type("html").send(req.query.html); }
  @Mutation(() => String)
  htmlControl(@Context("req") req: any, @Context("res") res: any) { return res.json({ value: req.query.html }); }
  @Mutation(() => String)
  htmlUncertain(@Context("req") req: any, @Context("res") res: any) { return res.type("html").send(fragment(req.query.html)); }
  @Mutation(() => String)
  csrfPositive(@Context("req") req: any, @Context("res") res: any) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
  @Mutation(() => String)
  csrfControl(@Context("req") req: any, @Context("res") res: any) { if (!tokens.verify(req.session.csrfSecret, req.headers["x-csrf-token"])) { throw new Error("invalid token"); } req.session.counter = req.body.counter; return res.send({ ok: true }); }
  @Mutation(() => String)
  csrfUncertain(@Context("req") req: any, @Context("res") res: any) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
  @Mutation(() => String)
  async authPositive(@Context("req") req: any, @Context("res") res: any) { if (req.user.role !== "writer") { res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  @Mutation(() => String)
  async authControl(@Context("req") req: any, @Context("res") res: any) { if (req.user.role !== "writer") { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  @Mutation(() => String)
  async authUncertain(@Context("req") req: any, @Context("res") res: any) { if (req.user.role !== "writer" && !enabled(req)) { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  @Mutation(() => String)
  redirectPositive(@Context("req") req: any, @Context("res") res: any) { return res.redirect(req.query.target); }
  @Mutation(() => String)
  redirectControl(@Context("req") req: any, @Context("res") res: any) { const target = req.query.target; if (!target.startsWith("/account/")) { return res.status(400).send("bad target"); } return res.redirect(target); }
  @Mutation(() => String)
  redirectUncertain(@Context("req") req: any, @Context("res") res: any) { return res.redirect(fragment(req.query.target)); }
  @Mutation(() => String)
  async massPositive(@Context("req") req: any, @Context("res") res: any) { return prisma.user.create({ data: { tenant: req.user.tenant, ...req.body } }); }
  @Mutation(() => String)
  async massControl(@Context("req") req: any, @Context("res") res: any) { return prisma.user.create({ data: { name: req.body.name, tenant: req.user.tenant } }); }
  @Mutation(() => String)
  async massUncertain(@Context("req") req: any, @Context("res") res: any) { return prisma.user.create({ data: { tenant: req.user.tenant, ...transform(req.body) } }); }
  @Mutation(() => String)
  corsPositive(@Context("req") req: any, @Context("res") res: any) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", req.headers.origin); return res.json({ ok: true }); }
  @Mutation(() => String)
  corsControl(@Context("req") req: any, @Context("res") res: any) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", "https://trusted.example"); return res.json({ ok: true }); }
  @Mutation(() => String)
  corsUncertain(@Context("req") req: any, @Context("res") res: any) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", policyOrigin(req.headers.origin)); return res.json({ ok: true }); }
  @Mutation(() => String)
  sessionUncertain(@Context("req") req: any, @Context("res") res: any) { return res.cookie("sid", req.session.id, {secure:enabled(req),httpOnly:true}).send({ok:true}); }
  @Mutation(() => String)
  tokenIssue(@Context("req") req: any, @Context("res") res: any) { req.session.csrfSecret=tokens.secretSync(); return res.json({token:tokens.create(req.session.csrfSecret)}); }
}
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
export class UpdateDto { readonly tenant!: string; name!: string; }
export class TypedController { @Post("typed") async typed(@Body() body:UpdateDto,@Req() req:any){return prisma.user.create({data:body});} }
Controller("typed")(TypedController);
export class MetadataController {
@ReadAccess(Modes.open)
@Mutation(() => String)
publicOperation(){return "observed";}
@ReadAccess(Modes.writer)
@Mutation(() => String)
requiredOperation(){return "observed";}
}
Resolver()(MetadataController);
export class MetadataModule{}
Module({providers:[MetadataController,{provide:APP_GUARD,useClass:SourceGuard}]})(MetadataModule);
