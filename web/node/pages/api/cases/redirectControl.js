import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import { fragment,transform,enabled,policyOrigin } from "../../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export default async function redirectControl(req, res) {const session=await getServerSession({req,res}); const target = req.query.target; if (!target.startsWith("/account/")) { return res.status(400).send("bad target"); } return res.redirect(target); }
