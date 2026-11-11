package fixture
import kotlin.collections.listOf
import kotlin.collections.setOf
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.PropertySource
import org.springframework.context.annotation.PropertySources
import org.springframework.beans.factory.annotation.Value
import org.springframework.session.web.http.*
import org.springframework.session.MapSessionRepository
import org.springframework.session.config.annotation.web.http.EnableSpringHttpSession
import org.springframework.web.cors.CorsConfiguration
import org.springframework.web.cors.UrlBasedCorsConfigurationSource
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.web.SecurityFilterChain
import org.springframework.security.authorization.*
import org.springframework.security.web.access.intercept.RequestAuthorizationContext
import org.thymeleaf.spring6.templateresolver.SpringResourceTemplateResolver
import org.thymeleaf.templatemode.TemplateMode
import java.util.function.UnaryOperator
@org.springframework.context.annotation.Configuration
@PropertySources(value = [PropertySource(value = ["classpath:cookie-settings.xml"]), PropertySource(value = ["classpath:cookie-settings.json"], factory = JsonSettingsFactory::class), PropertySource(value = ["classpath:cookie-settings.properties"], encoding = "UTF-8"), PropertySource(value = ["classpath:cookie-latin1.properties"])])
@EnableSpringHttpSession
open class Configuration {
  @Bean open fun policyChain(http: HttpSecurity, manager: AuthorizationManager<RequestAuthorizationContext>, source: UrlBasedCorsConfigurationSource): SecurityFilterChain {
    http.authorizeHttpRequests { registry -> registry.anyRequest().access(manager) }.cors { cors -> cors.configurationSource(source) }.csrf { csrf -> csrf.disable() }; return http.build()
  }
  @Bean open fun manager(): AuthorizationManager<RequestAuthorizationContext> { return AuthorizationManager { authentication, context -> AuthorizationDecision(authentication.get().isAuthenticated()) } }
  @Bean open fun sessions(@Value("\${web.cookie.secure}") secure: Boolean, @Value("\${web.cookié.http-only}") httpOnly: Boolean): CookieHttpSessionIdResolver {
    val serializer: DefaultCookieSerializer = DefaultCookieSerializer(); serializer.setUseSecureCookie(secure); serializer.setUseHttpOnlyCookie(httpOnly)
    val resolver: CookieHttpSessionIdResolver = CookieHttpSessionIdResolver(); resolver.setCookieSerializer(serializer); return resolver
  }
  @Bean open fun sessionRepository(): MapSessionRepository { return MapSessionRepository(java.util.concurrent.ConcurrentHashMap()) }
  @Bean open fun origins(@Value("\${web.origins}") origins: String): UrlBasedCorsConfigurationSource {
    val policy: CorsConfiguration = CorsConfiguration(); policy.setAllowCredentials(true); policy.setAllowedOriginPatterns(origins.split(","))
    val source: UrlBasedCorsConfigurationSource = UrlBasedCorsConfigurationSource(); source.registerCorsConfiguration("/**", policy); return source
  }
  @Bean open fun renderer(): UnaryOperator<String> { return UnaryOperator { value -> value.trim() } }
  @Bean open fun xmlResolver(): SpringResourceTemplateResolver { val resolver: SpringResourceTemplateResolver = SpringResourceTemplateResolver(); resolver.setPrefix("classpath:/templates/"); resolver.setSuffix(".xhtml"); resolver.setTemplateMode(TemplateMode.XML); resolver.setResolvablePatterns(setOf("xml-*")); resolver.setOrder(0); return resolver }
}
