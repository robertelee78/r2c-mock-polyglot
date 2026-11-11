package fixture
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.ui.Model
@Controller
class Views {
  @GetMapping("/view") fun page(model: Model, @RequestParam input: String): String { model.addAttribute("content", input); return "page" }
  @GetMapping("/xml") fun xmlPage(model: Model, @RequestParam input: String): String { model.addAttribute("content", input); return "xml-page" }
}
