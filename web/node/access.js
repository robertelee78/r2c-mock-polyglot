import {SetMetadata} from "@nestjs/common";
export const ReadAccess=(...values)=>SetMetadata("fixture:access",values);
