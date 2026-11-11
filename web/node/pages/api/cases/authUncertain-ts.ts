import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import {fragment,transform,enabled,policyOrigin} from "../../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export default async function authUncertain(req:any,res:any) {const session=await getServerSession({req,res}); if (session.user.role !== "writer" && !enabled(req)) { return res.status(403).send("denied"); } return prisma.user.update({ where: {id: session.user.id}, data: {name: req.body.name} }); }
