import { Entity, EntityType, IDManager } from './entity.js';

class Character extends Entity {
  getEntityType() {
    return EntityType.CHARACTER;
  }
}

// Create and register a new character
const idManager = IDManager.getInstance();
const character = new Character();
idManager.getAvailableIDNum(character);
idManager.register(character);

console.log(character.getID()); // Outputs the assigned ID
console.log(idManager.find(character.getID())); // Retrieves the character by ID