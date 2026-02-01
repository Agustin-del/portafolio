CREATE TABLE proyectos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    titulo TEXT NOT NULL UNIQUE,
    descripcion_general TEXT NOT NULL,
    descripcion_detallada TEXT NOT NULL
);

CREATE TABLE archivos_src (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id INTEGER NOT NULL,
    nombre TEXT NOT NULL,
    contenido TEXT NOT NULL,
    FOREIGN KEY(proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE
);
