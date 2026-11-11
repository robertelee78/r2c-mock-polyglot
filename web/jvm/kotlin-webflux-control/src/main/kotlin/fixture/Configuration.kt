package fixture
import kotlin.collections.listOf
import kotlin.collections.setOf
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.PropertySource
import org.springframework.context.annotation.PropertySources
import org.springframework.beans.factory.annotation.Value
import org.springframework.web.server.session.*
import org.springframework.web.cors.CorsConfiguration
import org.springframework.web.cors.reactive.UrlBasedCorsConfigurationSource
import org.springframework.security.config.web.server.ServerHttpSecurity
import org.springframework.security.web.server.SecurityWebFilterChain
import org.springframework.security.authorization.*
import org.springframework.security.web.server.authorization.AuthorizationContext
import org.thymeleaf.spring6.templateresolver.SpringResourceTemplateResolver
import org.thymeleaf.templatemode.TemplateMode
import java.util.function.UnaryOperator
@org.springframework.context.annotation.Configuration
@PropertySources(value = [PropertySource(value = ["classpath:cookie-settings.xml"]), PropertySource(value = ["classpath:cookie-settings.json"], factory = JsonSettingsFactory::class), PropertySource(value = ["classpath:cookie-settings.properties"], encoding = "UTF-8")])
open class Configuration {
  @Bean open fun policyChain(http: ServerHttpSecurity, manager: ReactiveAuthorizationManager<AuthorizationContext>, source: UrlBasedCorsConfigurationSource): SecurityWebFilterChain {
    http.authorizeExchange { registry -> registry.anyExchange().permitAll() }.cors { cors -> cors.configurationSource(source) }.csrf { csrf -> csrf.disable() }; return http.build()
  }
  @Bean open fun manager(): ReactiveAuthorizationManager<AuthorizationContext> { return ReactiveAuthorizationManager { authentication, context -> authentication.map { auth -> AuthorizationDecision(auth.isAuthenticated()) } } }
  @Bean(name = ["webSessionManager"]) open fun sessions(@Value("\${web.cookie.secure}") secure: Boolean): WebSessionManager {
    val resolver: CookieWebSessionIdResolver = CookieWebSessionIdResolver(); resolver.addCookieInitializer { cookie -> cookie.secure(secure).httpOnly(true) }
    val manager: DefaultWebSessionManager = DefaultWebSessionManager(); manager.setSessionIdResolver(resolver); return manager
  }
  @Bean open fun origins(@Value("\${web.origins}") origins: String): UrlBasedCorsConfigurationSource {
    val policy: CorsConfiguration = CorsConfiguration(); policy.setAllowCredentials(true); policy.setAllowedOriginPatterns(listOf("https://trusted.example"))
    val source: UrlBasedCorsConfigurationSource = UrlBasedCorsConfigurationSource(); source.registerCorsConfiguration("/**", policy); return source
  }
  @Bean open fun renderer(): UnaryOperator<String> { return UnaryOperator { value -> value.trim() } }
  @Bean open fun xmlResolver(): SpringResourceTemplateResolver { val resolver: SpringResourceTemplateResolver = SpringResourceTemplateResolver(); resolver.setPrefix("classpath:/templates/"); resolver.setSuffix(".xhtml"); resolver.setTemplateMode(TemplateMode.XML); resolver.setResolvablePatterns(setOf("xml-*")); resolver.setOrder(0); return resolver }
}
