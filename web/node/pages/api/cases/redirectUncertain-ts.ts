import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import { fragment,transform,enabled,policyOrigin } from "../../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
import type { NextApiRequest,NextApiResponse } from "next";
export default async function redirectUncertain(req:NextApiRequest & {user:any}, res:NextApiResponse) {const session=await getServerSession({req,res}); return res.redirect(fragment(req.query.target)); }
