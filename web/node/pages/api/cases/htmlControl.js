import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import { fragment,transform,enabled,policyOrigin } from "../../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export default async function htmlControl(req, res) {const session=await getServerSession({req,res}); return res.json({ value: req.query.html }); }
