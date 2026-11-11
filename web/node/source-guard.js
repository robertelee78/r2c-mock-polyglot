import {Reflector} from "@nestjs/core";
import {Modes} from "./permissions.js";
export class SourceGuard {
 constructor(){this.reflector=new Reflector();}
 canActivate(context){const required=this.reflector.get("fixture:access",context.getHandler());const bypass=!!required&&required.includes(Modes.open);if(bypass){return true;}return this.evaluate(required,context);}
 evaluate(required,context){const user=context.switchToHttp().getRequest().user;return !!user&&Array.isArray(required)&&required.every(permission=>user.permissions.includes(permission));}
}
