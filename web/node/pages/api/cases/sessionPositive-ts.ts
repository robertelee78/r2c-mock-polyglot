import NextAuth from "next-auth";
import {enabled} from "../../../helpers.js";
export default async function sessionPositive(req:any,res:any){return NextAuth(req,res,{providers:[],secret:"controlled-corpus-only",cookies:{sessionToken:{name:"sid",options:{secure:false,httpOnly:false}}}});}
