/**~~~includes: purpose, author, version, date, license, dependencies~~ 
 * module_template.js
 * Template for browser feature modules
 */

/** add a commit-hook that ensure only JS that passes JSLint can be submitted to a code repository */

/** global $, spa 
 *  example: spa.shell namespace, the filename would be spa.shell.js
 *  1. Create a namespace for the module using a self-executing function
 *  2. This prevents accidental creation of global JS variables
 *  3. Only 1 namespace should be defined per file
 *     and the filename should correlate precisely
*/
spa.module = (async() => {

  //---------------- BEGIN MODULE SCOPE VARIABLES ----------------
  /** 1. Declare and initialize module-scope variables.
   *  2. We commonly include 
   *     -configMap: to store module configurations
   *     -stateMap: to store run-time state values
   *     -cacheMap: to cache Query collections
   */
  var
    configMap = {
      settable_map: { color_name: true },
      color_name: '#8e8'
    }
    // stateMap = { target_push: undefined },
    // stateSet = {  },
    
    // setMap, configModule, initModule;
  //---------------- END MODULE SCOPE VARIABLES ----------------
  //------------------ BEGIN UTILITY METHODS -------------------
  // example: getTrimmedString
  /** 1. Group ALL private util methods in their own section
   *  2. These methods don't manipulate the DOM and therefore
   *     don't require a browser to run
   *  3. If a method has a utility beyond a single module
   *     -> put it in the shared utils section
   */
  //------------------- END UTILITY METHODS --------------------
  //-------------------- BEGIN DOM METHODS ---------------------
  /** 1. Generate all private DOM methods in their own section
   *  2. These methods access and modify the DOM and therefore
   *     require a browser to run
   *  3. An example DOM method: move_css_sprite
   */
  //------------------- END DOM METHODS ------------------------
}) ();
