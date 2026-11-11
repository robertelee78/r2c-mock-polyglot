"use server";
import { PrismaClient } from "@prisma/client";
import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { cookies } from "next/headers";
import Tokens from "csrf";
import {redirect} from "next/navigation";
import {fragment,transform,enabled} from "../helpers.js";
const prisma=new PrismaClient();const tokens=new Tokens();
export async function redirectPositive(input){return redirect(input.target);}
export async function massPositive(input){const session=await getServerSession();return prisma.user.create({data:{tenant:session.user.tenant,...input}});}
export async function sessionPositive(input){const session=await getServerSession();cookies().set("sid",session.user.id,{secure:false,httpOnly:false});return {ok:true};}
export async function csrfPositive(input){const session=await getServerSession();return prisma.user.update({where:{id:session.user.id},data:{name:input.name}});}
export async function authPositive(input){const session=await getServerSession();if(session.user.role!=="writer"){NextResponse.json({denied:true},{status:403});}return prisma.user.update({where:{id:session.user.id},data:{name:input.name}});}
export async function redirectControl(input){const target=input.target;if(!target.startsWith("/account/")){throw new Error("bad target");}return redirect(target);}
export async function massControl(input){const session=await getServerSession();return prisma.user.create({data:{name:input.name,tenant:session.user.tenant}});}
export async function sessionControl(input){const session=await getServerSession();cookies().set("sid",session.user.id,{secure:true,httpOnly:true});return {ok:true};}
export async function csrfControl(input){const session=await getServerSession();if(!tokens.verify(session.csrfSecret,input.csrfToken)){throw new Error("invalid token");}return prisma.user.update({where:{id:session.user.id},data:{name:input.name}});}
export async function authControl(input){const session=await getServerSession();if(session.user.role!=="writer"){throw new Error("denied");}return prisma.user.update({where:{id:session.user.id},data:{name:input.name}});}
export async function redirectUncertain(input){return redirect(fragment(input.target));}
export async function massUncertain(input){const session=await getServerSession();return prisma.user.create({data:{tenant:session.user.tenant,...transform(input)}});}
export async function sessionUncertain(input){const session=await getServerSession();cookies().set("sid",session.user.id,{secure:enabled(input),httpOnly:true});return {ok:true};}
export async function csrfUncertain(input){const session=await getServerSession();if(!enabled(input)){throw new Error("blocked");}return prisma.user.update({where:{id:session.user.id},data:{name:input.name}});}
export async function authUncertain(input){const session=await getServerSession();if(session.user.role!=="writer"&&!enabled(input)){throw new Error("denied");}return prisma.user.update({where:{id:session.user.id},data:{name:input.name}});}
