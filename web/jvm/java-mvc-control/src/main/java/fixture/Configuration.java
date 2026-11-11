package fixture;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.PropertySource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.session.web.http.*;
import org.springframework.session.MapSessionRepository;
import org.springframework.session.config.annotation.web.http.EnableSpringHttpSession;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.authorization.*;
import org.springframework.security.web.access.intercept.RequestAuthorizationContext;
import org.thymeleaf.spring6.templateresolver.SpringResourceTemplateResolver;
import org.thymeleaf.templatemode.TemplateMode;
import java.util.function.UnaryOperator;
@org.springframework.context.annotation.Configuration
@PropertySource(value = "classpath:cookie-settings.properties", encoding = "UTF-8")
@PropertySource("classpath:cookie-settings.xml")
@PropertySource(value = "classpath:cookie-settings.json", factory = JsonSettingsFactory.class)
@EnableSpringHttpSession
public class Configuration {
  @Bean public SecurityFilterChain policyChain(HttpSecurity http, AuthorizationManager<RequestAuthorizationContext> manager, UrlBasedCorsConfigurationSource source) throws Exception {
    http.authorizeHttpRequests(registry -> registry.anyRequest().permitAll()).cors(cors -> cors.configurationSource(source)).csrf(csrf -> csrf.disable()); return http.build();
  }
  @Bean public AuthorizationManager<RequestAuthorizationContext> manager() { return (authentication, context) -> new AuthorizationDecision(authentication.get().isAuthenticated()); }
  @Bean public CookieHttpSessionIdResolver sessions(@Value("${web.cookie.secure}") boolean secure) {
    DefaultCookieSerializer serializer = new DefaultCookieSerializer(); serializer.setUseSecureCookie(secure); serializer.setUseHttpOnlyCookie(true);
    CookieHttpSessionIdResolver resolver = new CookieHttpSessionIdResolver(); resolver.setCookieSerializer(serializer); return resolver;
  }
  @Bean public MapSessionRepository sessionRepository() { return new MapSessionRepository(new java.util.concurrent.ConcurrentHashMap<>()); }
  @Bean public UrlBasedCorsConfigurationSource origins(@Value("${web.origins}") String origins) {
    CorsConfiguration policy = new CorsConfiguration(); policy.setAllowCredentials(true); policy.setAllowedOriginPatterns(java.util.List.of("https://trusted.example"));
    UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource(); source.registerCorsConfiguration("/**", policy); return source;
  }
  @Bean public UnaryOperator<String> renderer() { return value -> value.trim(); }
  @Bean public SpringResourceTemplateResolver xmlResolver() { SpringResourceTemplateResolver resolver = new SpringResourceTemplateResolver(); resolver.setPrefix("classpath:/templates/"); resolver.setSuffix(".xhtml"); resolver.setTemplateMode(TemplateMode.XML); resolver.setResolvablePatterns(java.util.Set.of("xml-*")); resolver.setOrder(0); return resolver; }
}
