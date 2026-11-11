import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import { fragment,transform,enabled,policyOrigin } from "../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export async function POST(req){const session = await getServerSession(); const body = await req.json(); return NextResponse.json(await prisma.user.create({data:{tenant:session.user.tenant,name:body.name}}));}
