package fixture
import org.springframework.boot.json.JsonParserFactory
import org.springframework.core.env.MapPropertySource
import org.springframework.core.env.PropertySource
import org.springframework.core.io.support.EncodedResource
import org.springframework.core.io.support.PropertySourceFactory
class JsonSettingsFactory : PropertySourceFactory {
  override fun createPropertySource(name: String?, resource: EncodedResource): PropertySource<*> {
    val json = resource.resource.inputStream.use { stream -> String(stream.readBytes(), Charsets.UTF_8) }
    return MapPropertySource(name ?: "fixture-json", JsonParserFactory.getJsonParser().parseMap(json))
  }
}
