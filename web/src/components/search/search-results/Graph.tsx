/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
import '@react-sigma/core/lib/react-sigma.min.css';

import { Edge, Node } from '../../../shared/models/SearchModel';
import { SigmaContainer, useLoadGraph } from '@react-sigma/core';

import MultiDirectedGraph from 'graphology';
import React from 'react';
import { useWorkerLayoutForceAtlas2 } from '@react-sigma/layout-forceatlas2';

type LoadGraphProps = {
  nodes: Node[];
  edges: Edge[];
};

export const LoadGraph = (props: LoadGraphProps) => {
  const { nodes, edges } = props;
  const load = useLoadGraph();

  //eslint-disable-next-line
  const layout = useWorkerLayoutForceAtlas2({
    settings: {
      gravity: 1,
    },
  });

  //load and animate graph
  React.useEffect(() => {
    const graph = new MultiDirectedGraph();

    nodes.forEach((node) => {
      graph.mergeNode(node.oid, {
        label: node.name,
        color: '#000',
        size: 10,
        x: Math.random(),
        y: Math.random(),
      });
    });

    edges.forEach((edge) => {
      let label = '';
      edge.subjectNames.forEach((subjectName) => {
        //get lat part of subjectName after last slash
        const lat = subjectName.substring(subjectName.lastIndexOf('/') + 1);
        label = label + lat + '    ';
      });

      graph.addEdgeWithKey(edge.from_oid + edge.to_oid, edge.from_oid, edge.to_oid, {
        label: label,
        color: '#000',
        edgeType: 'arrow',
      });
    });

    load(graph);
    layout.start();

    layout.start();
    return layout.stop;
  }, [nodes, edges, load, layout]);

  return null;
};

export const DisplayGraph = (props: LoadGraphProps) => {
  if (!props.nodes || !props.edges) {
    return null;
  }

  return (
    <SigmaContainer
      style={{ position: 'absolute', top: 150, left: 0, zIndex: 1, height: '80vh' }}
      graph={MultiDirectedGraph}
      settings={{ renderEdgeLabels: true }}
    >
      <LoadGraph nodes={props.nodes} edges={props.edges} />
    </SigmaContainer>
  );
};

//a84aa16f5ffd26a40792268a5febf8e8ff468db1
