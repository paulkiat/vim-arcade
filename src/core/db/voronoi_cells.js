async function createVoronoiCells() {
  return new Promise((resolve, reject) => {
    resolve(voronoi.compute(points));
  });

  await createVoronoiCells();
}