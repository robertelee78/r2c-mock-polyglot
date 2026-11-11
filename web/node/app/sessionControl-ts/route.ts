import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import { NextRequest } from "next/server";
import { fragment,transform,enabled,policyOrigin } from "../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export async function POST(req: NextRequest){const session=await getServerSession();cookies().set("sid",session.user.id,{secure:true,httpOnly:true});return NextResponse.json({ok:true});}
