import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import { fragment,transform,enabled,policyOrigin } from "../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export async function POST(req){const session = await getServerSession(); if (session.user.role !== "writer" && !enabled(req)) { return NextResponse.json({denied:true},{status:403}); } return NextResponse.json(await prisma.user.update({where:{id:session.user.id},data:{name:(await req.json()).name}}));}
