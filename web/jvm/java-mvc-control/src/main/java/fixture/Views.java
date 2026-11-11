package fixture;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.ui.Model;
@Controller
public class Views {
  @GetMapping("/view") public String page(Model model, @RequestParam String input) { model.addAttribute("content", input); return "page"; }
  @GetMapping("/xml") public String xmlPage(Model model, @RequestParam String input) { model.addAttribute("content", input); return "xml-page"; }
}
