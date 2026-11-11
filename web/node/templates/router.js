import express from "express";
import {mountedRaw,mountedEscaped,mountedPartial} from "./handlers.js";
export const router=express.Router();
router.get("/raw",mountedRaw);
router.get("/escaped",mountedEscaped);
router.get("/partial",mountedPartial);
