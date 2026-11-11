package fixture;
import jakarta.persistence.*;
import com.fasterxml.jackson.annotation.JsonProperty;
@Entity
public class Account {
  @Id @GeneratedValue private Long id;
  @JsonProperty(access = JsonProperty.Access.READ_ONLY)
  private String owner;
  private String displayName;
  public Long getId() { return id; }
  public void setId(Long value) { id = value; }
  public String getOwner() { return owner; }
  public void setOwner(String value) { owner = value; }
  public String getDisplayName() { return displayName; }
  public void setDisplayName(String value) { displayName = value; }
}
