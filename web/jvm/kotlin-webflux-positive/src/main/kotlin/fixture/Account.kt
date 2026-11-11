package fixture
import jakarta.persistence.*
import com.fasterxml.jackson.annotation.JsonProperty
@Entity
open class Account {
  @Id @GeneratedValue private var id: Long? = null
  @field:JsonProperty(access = JsonProperty.Access.READ_ONLY)
  open var owner: String? = null
  private var displayNameValue: String? = null
  open fun getId(): Long? { return id }
  open fun setId(value: Long?) { id = value }
  open fun getDisplayName(): String? { return displayNameValue }
  open fun setDisplayName(value: String?) { displayNameValue = value }
}
