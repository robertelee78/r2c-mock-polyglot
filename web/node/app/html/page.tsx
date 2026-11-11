import React from "react";
export function HtmlPositive({value}: {value:string}) { return <article dangerouslySetInnerHTML={{__html:value}} />; }
export function HtmlControl({value}: {value:string}) { return <article>{value}</article>; }
export function HtmlUncertain({value, render}: {value:string,render:(value:string)=>string}) { return <article dangerouslySetInnerHTML={{__html:render(value)}} />; }
export default HtmlControl;
