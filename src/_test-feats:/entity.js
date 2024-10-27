// Use an IIFE to create a module and avoid polluting the global scope
const EntityModule = (function() {
  // Private variables and functions
  const EntityType = {
    UNDEFINED: -1,
    CHARACTER: 0,
    OBJECT: 1,
    ROOM: 2,
    QA_ENTRY: 3
  };
  
  // Prototype methods for Entity
  Entity.prototype = {
    getID: function() { return this.id; },
    getScript: function() { return this.script; },
    isPurged: function() { return this.purged; },
    getIDString: function() { return `need a return val`}, // todo fin later
    extract: function() { /* Implementation */ },
    setID: function(id) { this.id = id; },
    setPurged: function() { this.purged = true; },
    purge: function() { /* Implementation */ },
    getEntityType: function() { throw new Error("getEntityType must be implemented by subclasses"); }
  };

  // Singleton IDManager
  const IDManager = (function() {
    let instance;

    function init() {
      let entities = new Map();
      let nextID = 1;

      return {
        getAvailableIDNum: function(entity) {
          const id = nextID++;
          entity.setID(id);
          return id;
        },
        register: function(entity) {
          entities.set(entity.getID(), entity);
        },
        unregister: function(entity) {
          entities.delete(entity.getID());
        },
        find: function(id) {
          return entities.get(id) || null;
        },
        getAllEntities: function() {
          return Array.from(entities.values());
        }
      };
    }

    return {
      getInstance: function() {
        if (!instance) {
          instance = init();
        }
        return instance;
      }
    };
  })();

  // QAEntry constructor
  function QAEntry(question, answer) {
    Entity.call(this);
    this.question = question;
    this.answer = answer;
  }

  // Inherit from Entity
  QAEntry.prototype = Object.create(Entity.prototype);
  QAEntry.prototype.constructor = QAEntry;
  QAEntry.prototype.getEntityType = function() { return EntityType.QA_ENTRY; };

  // Public API
  return {
    Entity: Entity,
    EntityType: EntityType,
    IDManager: IDManager,
    QAEntry: QAEntry
  };
})();

export default EntityModule;