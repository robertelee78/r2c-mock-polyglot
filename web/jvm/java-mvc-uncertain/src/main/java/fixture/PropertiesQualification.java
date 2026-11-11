package fixture;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.PropertySource;
import org.springframework.session.web.http.CookieHttpSessionIdResolver;
import org.springframework.session.web.http.DefaultCookieSerializer;
@org.springframework.context.annotation.Configuration
@PropertySource("classpath:qualification-duplicate.properties")
@PropertySource(value = "classpath:qualification-malformed.properties", encoding = "UTF-8")
@PropertySource(value = "classpath:qualification-charset.properties", encoding = "windows-1252")
public class PropertiesQualification {
  @Bean public CookieHttpSessionIdResolver duplicateSessions(@Value("${qualification.duplicate.secure}") boolean secure) {
    DefaultCookieSerializer serializer = new DefaultCookieSerializer(); serializer.setUseSecureCookie(secure);
    CookieHttpSessionIdResolver resolver = new CookieHttpSessionIdResolver(); resolver.setCookieSerializer(serializer); return resolver;
  }
  @Bean public CookieHttpSessionIdResolver malformedSessions(@Value("${qualification.malformed.secure}") boolean secure, @Value("${qualification.malformed.http-only}") boolean httpOnly) {
    DefaultCookieSerializer serializer = new DefaultCookieSerializer(); serializer.setUseSecureCookie(secure); serializer.setUseHttpOnlyCookie(httpOnly);
    CookieHttpSessionIdResolver resolver = new CookieHttpSessionIdResolver(); resolver.setCookieSerializer(serializer); return resolver;
  }
  @Bean public CookieHttpSessionIdResolver charsetSessions(@Value("${qualification.charset.secure}") boolean secure) {
    DefaultCookieSerializer serializer = new DefaultCookieSerializer(); serializer.setUseSecureCookie(secure);
    CookieHttpSessionIdResolver resolver = new CookieHttpSessionIdResolver(); resolver.setCookieSerializer(serializer); return resolver;
  }
}
