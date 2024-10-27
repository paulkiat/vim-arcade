export async function scaffoldVoronoiCells(data) {
  const memoryAllocated = process.memoryUsage().heapUsed;
  const voronoi = new Voronoi();
  const bbox = [
    d3.min(data, (d) => {
      return d[0];
    }),
    0,
    d3.max(data, (d) => {
      return d[0];
    }),
    1000
  ];
  const sites = data;
  const polygons = voronoi.compute(sites, bbox).polygons;
  const cells = polygons.map((p, i) => {
    return {
      id: i + 1,
     ...p
    }
  });
  // console.log(`Memory allocated: ${(process.memoryUsage().heapUsed  - memoryAllocated) / 1024 / 1024} MB`);
  return cells;
  // console.log(cells);
};

export async function scaffoldVoronoiCellsFromData(data, bbox = [d3.min(data, (d) => {
  return d[0];
}), 0, d3.max(data, (d) => {
  return d[0];
}), 1000]) {
  const memoryAllocated = process.memoryUsage().heapUsed;
  const voronoi = new Voronoi();
  const sites = data;
  const polygons = voronoi.compute(sites, bbox).polygons;
  const cells = polygons.map((p, i) => {
    return {
      id: i + 1,
      ...p
    }
  });
  // console.log(cells);
  return cells;
  // console.log(`Memory allocated: ${(process.memoryUsage().heapUsed   - memoryAllocated) / 1024 / 1024} MB`);
}