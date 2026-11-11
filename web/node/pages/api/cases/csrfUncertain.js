import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import {enabled} from "../../../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export default async function csrfUncertain(req,res){const session=await getServerSession({req,res});if(!enabled(req)){return res.status(403).json({denied:true});}return res.json(await prisma.user.update({where:{id:session.user.id},data:{name:req.body.name}}));}
