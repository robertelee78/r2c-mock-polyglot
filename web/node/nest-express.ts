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
@Post("sessionControl")
sessionControl(@Req() req:any,@Res() res:any) { return res.cookie("sid",req.session.id,{secure:true,httpOnly:true}).send({ok:true}); }
@Post("sessionPositive")
sessionPositive(@Req() req:any,@Res() res:any) { return res.cookie("sid",req.session.id,{secure:false,httpOnly:false}).send({ok:true}); }
  @Post("htmlPositive")
  htmlPositive(@Req() req: any, @Res() res: any) { return res.type("html").send(req.query.html); }
  @Post("htmlControl")
  htmlControl(@Req() req: any, @Res() res: any) { return res.json({ value: req.query.html }); }
  @Post("htmlUncertain")
  htmlUncertain(@Req() req: any, @Res() res: any) { return res.type("html").send(fragment(req.query.html)); }
  @Post("csrfPositive")
  csrfPositive(@Req() req: any, @Res() res: any) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
  @Post("csrfControl")
  csrfControl(@Req() req: any, @Res() res: any) { if (!tokens.verify(req.session.csrfSecret, req.headers["x-csrf-token"])) { throw new Error("invalid token"); } req.session.counter = req.body.counter; return res.send({ ok: true }); }
  @Post("csrfUncertain")
  csrfUncertain(@Req() req: any, @Res() res: any) { req.session.counter = req.body.counter; return res.json({ ok: true }); }
  @Post("authPositive")
  async authPositive(@Req() req: any, @Res() res: any) { if (req.user.role !== "writer") { res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  @Post("authControl")
  async authControl(@Req() req: any, @Res() res: any) { if (req.user.role !== "writer") { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  @Post("authUncertain")
  async authUncertain(@Req() req: any, @Res() res: any) { if (req.user.role !== "writer" && !enabled(req)) { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: req.user.id}, data: {name: req.body.name} }); }
  @Post("redirectPositive")
  redirectPositive(@Req() req: any, @Res() res: any) { return res.redirect(req.query.target); }
  @Post("redirectControl")
  redirectControl(@Req() req: any, @Res() res: any) { const target = req.query.target; if (!target.startsWith("/account/")) { return res.status(400).send("bad target"); } return res.redirect(target); }
  @Post("redirectUncertain")
  redirectUncertain(@Req() req: any, @Res() res: any) { return res.redirect(fragment(req.query.target)); }
  @Post("massPositive")
  async massPositive(@Req() req: any, @Res() res: any) { return prisma.user.create({ data: { tenant: req.user.tenant, ...req.body } }); }
  @Post("massControl")
  async massControl(@Req() req: any, @Res() res: any) { return prisma.user.create({ data: { name: req.body.name, tenant: req.user.tenant } }); }
  @Post("massUncertain")
  async massUncertain(@Req() req: any, @Res() res: any) { return prisma.user.create({ data: { tenant: req.user.tenant, ...transform(req.body) } }); }
  @Post("corsPositive")
  corsPositive(@Req() req: any, @Res() res: any) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", req.headers.origin); return res.json({ ok: true }); }
  @Post("corsControl")
  corsControl(@Req() req: any, @Res() res: any) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", "https://trusted.example"); return res.json({ ok: true }); }
  @Post("corsUncertain")
  corsUncertain(@Req() req: any, @Res() res: any) { res.setHeader("Access-Control-Allow-Credentials", "true"); res.setHeader("Access-Control-Allow-Origin", policyOrigin(req.headers.origin)); return res.json({ ok: true }); }
  @Post("sessionUncertain")
  sessionUncertain(@Req() req: any, @Res() res: any) { return res.cookie("sid", req.session.id, {secure:enabled(req),httpOnly:true}).send({ok:true}); }
  @Post("tokenIssue")
  tokenIssue(@Req() req: any, @Res() res: any) { req.session.csrfSecret=tokens.secretSync(); return res.json({token:tokens.create(req.session.csrfSecret)}); }
}
Controller("corpus")(CorpusController);
export class CorpusModule {}
Module({controllers:[CorpusController]})(CorpusModule);
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
@Post("publicOperation")
publicOperation(){return "observed";}
@ReadAccess(Modes.writer)
@Post("requiredOperation")
requiredOperation(){return "observed";}
}
Controller("metadata")(MetadataController);
export class MetadataModule{}
Module({controllers:[MetadataController],providers:[{provide:APP_GUARD,useClass:SourceGuard}]})(MetadataModule);
