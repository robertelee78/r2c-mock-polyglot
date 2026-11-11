package fixture
class InputAccount {
  private var ownerValue: String? = null
  private var displayNameValue: String? = null
  fun getOwner(): String? { return ownerValue }
  fun setOwner(value: String?) { ownerValue = value }
  fun getDisplayName(): String? { return displayNameValue }
  fun setDisplayName(value: String?) { displayNameValue = value }
}
