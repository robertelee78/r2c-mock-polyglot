package fixture;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import org.springframework.boot.json.JsonParserFactory;
import org.springframework.core.env.MapPropertySource;
import org.springframework.core.env.PropertySource;
import org.springframework.core.io.support.EncodedResource;
import org.springframework.core.io.support.PropertySourceFactory;
public class JsonSettingsFactory implements PropertySourceFactory {
  @Override public PropertySource<?> createPropertySource(String name, EncodedResource resource) throws IOException {
    try (InputStream stream = resource.getResource().getInputStream()) {
      String json = new String(stream.readAllBytes(), StandardCharsets.UTF_8);
      return new MapPropertySource(name == null ? "fixture-json" : name, JsonParserFactory.getJsonParser().parseMap(json));
    }
  }
}
