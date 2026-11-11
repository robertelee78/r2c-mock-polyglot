import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import {fragment,transform,enabled,policyOrigin} from "../../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export default async function massUncertain(req:any,res:any) {const session=await getServerSession({req,res}); return prisma.user.create({ data: { tenant: session.user.tenant, ...transform(req.body) } }); }
