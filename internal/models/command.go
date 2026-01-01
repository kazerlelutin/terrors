package models

// Command représente une commande disponible dans le système CQRS
type Command struct {
	Name        string                 `json:"name"`              // Nom de la commande
	Description string                 `json:"description"`       // Description
	Method      string                 `json:"method"`            // HTTP method
	Path        string                 `json:"path"`              // Chemin de l'endpoint
	Params      map[string]interface{} `json:"params,omitempty"`  // Paramètres optionnels
	Example     map[string]interface{} `json:"example,omitempty"` // Exemple de requête
}

// DashboardResponse représente la réponse du dashboard avec toutes les commandes disponibles
type DashboardResponse struct {
	Message  string    `json:"message"`
	Quote    string    `json:"quote"`
	Commands []Command `json:"commands"`
}
