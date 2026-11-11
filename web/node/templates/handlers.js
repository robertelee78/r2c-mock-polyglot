export function mountedRaw(req,res){return res.render("raw.ejs",{value:req.query.value});}
export function mountedEscaped(req,res){return res.render("escaped.ejs",{value:req.query.value});}
export function mountedPartial(req,res){return res.render("include.hbs",{value:req.query.value});}
