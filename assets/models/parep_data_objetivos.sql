-- PAREP Objectives Data - Insert statements for all 7 objectives with results and indicators

-- PRÉ-ESCOLAR Objectives (ambito_id = 1)

-- Objetivo 1
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id) 
VALUES ('Assegurar que todas as crianças com idade de 4 e 5 anos frequentem o Pré-escolar', 
        'Acesso universal e inclusivo à Educação Pré-escolar com redução de assimetrias regionais e maior inclusão de crianças de famílias vulneráveis', 1);

-- Objetivo 1 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id) 
VALUES 
  ('Acesso universal e inclusivo à Educação Pré-escolar', 1),
  ('Redução das assimetrias regionais de acesso', 1),
  ('Aumento da frequência das crianças de 4 e 5 anos nos jardins de infância', 1),
  ('Maior inclusão de crianças provenientes de famílias vulneráveis', 1);

-- Objetivo 1 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND1', 'Taxa líquida de escolarização/cobertura das crianças de 4-5 anos no Pré-escolar', 84.1, 1, 1),
  ('IND2', 'Nº de crianças que frequentam Jardim de Infância', 15906, 2, 1),
  ('IND3', 'Nº de famílias de classe 1 e 2 que recebem apoio para escolarização na Educação Pré-escolar', NULL, 2, 1),
  ('IND4', 'Nº de Jardins que recebem Kits de materiais lúdico-pedagógicos e equipamentos para crianças com NEE', NULL, 2, 1);

-- Objetivo 1 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND1 2026', 95, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND1')),
  ('Meta IND2 2026', 17280, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND2')),
  ('Meta IND3 2026', 3500, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND3')),
  ('Meta IND4 2026', 150, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND4'));

---

-- Objetivo 2
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id)
VALUES ('Melhorar o desempenho das aprendizagens das crianças no Pré-escolar',
        'Melhoria das competências em leitura, escrita e numeracia com profissionais mais qualificados e contextos pedagógicos melhorados', 1);

-- Objetivo 2 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id)
VALUES 
  ('Melhoria das competências em leitura, escrita e numeracia', 2),
  ('Profissionais de infância mais qualificados', 2),
  ('Contextos pedagógicos de aprendizagem melhorados', 2),
  ('Práticas pedagógicas alinhadas com as orientações curriculares', 2);

-- Objetivo 2 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND5', '% de crianças de 4-5 que desenvolvem competências básicas em língua portuguesa e números', NULL, 1, 2),
  ('IND6', '% de Jardins que utiliza e cumpre o programa previsto em articulação com unidades educativas', NULL, 1, 2),
  ('IND7', 'Nº de educadoras com formação inicial formadas e integradas no sistema', 136, 2, 2),
  ('IND8', 'Nº de monitoras com formação inicial formadas e integradas no sistema', 283, 2, 2),
  ('IND9', '% de profissionais de infância com formação adequada', 9.6, 1, 2);

-- Objetivo 2 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND5 2026', 75, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND5')),
  ('Meta IND6 2026', 80, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND6')),
  ('Meta IND7 2026', 424, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND7')),
  ('Meta IND8 2026', 566, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND8')),
  ('Meta IND9 2026', 30, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND9'));

---

-- Objetivo 3
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id)
VALUES ('Melhorar a eficiência e eficácia do uso dos recursos disponibilizados ao ensino pré-escolar',
        'Planeamento eficiente da rede pré-escolar com melhor utilização de recursos humanos e infraestruturas', 1);

-- Objetivo 3 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id)
VALUES 
  ('Planeamento mais eficiente da rede pré-escolar', 3),
  ('Melhor utilização dos recursos humanos e infraestruturas', 3),
  ('Melhor gestão dos recursos disponibilizados ao subsistema', 3);

-- Objetivo 3 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND10', 'Carta Escolar do Pré-escolar elaborada', NULL, 2, 3),
  ('IND11', 'Rácio criança/profissional de infância', 20, 2, 3);

-- Objetivo 3 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND10 2026', 1, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND10')),
  ('Meta IND11 2026', 25, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND11'));

---

-- Objetivo 4
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id)
VALUES ('Reforçar a capacidade institucional e organizativa da Educação Pré-escolar',
        'Novo regime jurídico implementado, estatuto de carreira operacionalizado e gestão baseada em resultados', 1);

-- Objetivo 4 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id)
VALUES 
  ('Novo regime jurídico da Educação Pré-escolar implementado', 4),
  ('Estatuto de carreira dos profissionais de infância definido e operacionalizado', 4),
  ('Gestão dos jardins de infância baseada em resultados', 4),
  ('Maior articulação entre jardins de infância e escolas básicas', 4),
  ('Utilização do SIGE para monitoria e gestão do subsistema', 4);

-- Objetivo 4 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND12', '% de Jardins que cumprem os requisitos mínimos exigidos para funcionamento', NULL, 1, 4),
  ('IND13', 'Nº de Gestores de Jardins de Infância capacitados em gestão baseada em resultados', 582, 2, 4),
  ('IND14', '% de gestores dos Jardins habilitados para a função de direção', NULL, 1, 4),
  ('IND15', 'Nº de Jardins que fazem intercâmbio e articulam com as escolas básicas', NULL, 2, 4),
  ('IND16', 'Estatuto de carreira dos profissionais da infância elaborado e apropriado', NULL, 2, 4),
  ('IND17', 'Nº de Jardins com prática de gestão baseada em resultados', NULL, 2, 4),
  ('IND18', '% de Jardins com prática de gestão baseada em resultados', NULL, 1, 4),
  ('IND19', 'Nº de Gestores de jardins capacitados para uso e manuseamento do SIGE', 582, 2, 4),
  ('IND20', '% de Jardins que disponibiliza dados através do SIGE', NULL, 1, 4);

-- Objetivo 4 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND12 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND12')),
  ('Meta IND13 2026', 600, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND13')),
  ('Meta IND14 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND14')),
  ('Meta IND15 2026', 360, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND15')),
  ('Meta IND16 2026', 1, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND16')),
  ('Meta IND17 2026', 600, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND17')),
  ('Meta IND18 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND18')),
  ('Meta IND19 2026', 600, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND19')),
  ('Meta IND20 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND20'));

---

-- 1º CICLO Objectives (ambito_id = 2)

-- Objetivo 5
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id)
VALUES ('Consolidar o acesso equitativo e inclusivo no EBO',
        'Acesso universal e inclusivo consolidado com melhor integração de crianças com NEE e redução de desigualdades', 2);

-- Objetivo 5 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id)
VALUES 
  ('Acesso universal e inclusivo consolidado no EBO', 5),
  ('Melhor integração das crianças com Necessidades Educativas Especiais (NEE)', 5),
  ('Redução de desigualdades de género e vulnerabilidade social', 5),
  ('Reforço da capacidade de resposta das estruturas de educação inclusiva', 5);

-- Objetivo 5 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND21', 'Nº de campanhas regulares de sensibilização implementadas', NULL, 2, 5),
  ('IND22', 'Estudo de mapeamento de crianças com NEE elaborado', NULL, 2, 5),
  ('IND23', 'Equipas EMAEI reforçadas com mais um técnico em todos os concelhos', NULL, 2, 5),
  ('IND24', 'Reforço dos materiais lúdico-pedagógicos destinados às crianças com NEE', NULL, 2, 5),
  ('IND25', 'Taxa líquida de escolarização das crianças de 6-13 anos no EBO', 99.6, 1, 5),
  ('IND26', 'Índice de Paridade de Género (M/F)', 0.91, 3, 5);

-- Objetivo 5 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND21 2026', 20, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND21')),
  ('Meta IND22 2026', 1, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND22')),
  ('Meta IND23 2026', 22, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND23')),
  ('Meta IND24 2026', 22, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND24')),
  ('Meta IND25 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND25')),
  ('Meta IND26 2026', 0.95, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND26'));

---

-- Objetivo 6
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id)
VALUES ('Reforçar o êxito e a qualidade das aprendizagens em língua portuguesa e matemática',
        'Melhoria do desempenho em Português e Matemática com professores mais capacitados e maior inclusão de crianças com NEE', 2);

-- Objetivo 6 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id)
VALUES 
  ('Melhoria do desempenho dos alunos em Língua Portuguesa', 6),
  ('Melhoria do desempenho dos alunos em Matemática', 6),
  ('Professores mais capacitados na gestão curricular e avaliação', 6),
  ('Maior inclusão das crianças com NEE no processo de ensino-aprendizagem', 6),
  ('Aumento das taxas de sucesso escolar', 6);

-- Objetivo 6 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND27', 'Nº de docentes do 1º ciclo que recebem ações de capacitação em gestão curricular', 4188, 2, 6),
  ('IND28', 'Nº de docentes capacitados em gestão pedagógica de crianças com NEE', NULL, 2, 6),
  ('IND29', 'Nº de equipas pedagógicas capacitadas para supervisão e avaliação no EBO', 60, 2, 6),
  ('IND30', 'Percentagem de aprovação no 1º Ciclo do EBO', 83.9, 1, 6),
  ('IND31', 'Percentagem de reprovação no 1º Ciclo do EBO', 15.2, 1, 6),
  ('IND32', 'Percentagem de abandono no 1º Ciclo do EBO', 0.01, 1, 6),
  ('IND33', '% de crianças do EBO que adquirem competências básicas em LP e Matemática', 32.7, 1, 6),
  ('IND34', '% de docentes abrangidos pelo programa de formação contínua', NULL, 1, 6);

-- Objetivo 6 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND27 2026', 4200, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND27')),
  ('Meta IND28 2026', 4200, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND28')),
  ('Meta IND29 2026', 60, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND29')),
  ('Meta IND30 2026', 92, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND30')),
  ('Meta IND31 2026', 8, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND31')),
  ('Meta IND32 2026', 0, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND32')),
  ('Meta IND33 2026', 60, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND33')),
  ('Meta IND34 2026', 90, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND34'));

---

-- Objetivo 7
INSERT INTO objetivos (objetivo, desc_objetivo, ambito_id)
VALUES ('Melhorar a eficiência e eficácia do uso dos recursos disponibilizados no EBO',
        'Unidades educativas com gestão plenamente funcional, projetos educativos implementados e SIGE operacional', 2);

-- Objetivo 7 - Resultados Esperados
INSERT INTO resultado_esperados (resultado_esperado, objetivo_id)
VALUES 
  ('Unidades educativas com órgãos de gestão plenamente funcionais', 7),
  ('Generalização dos Projetos Educativos como instrumento de gestão', 7),
  ('Supervisão e acompanhamento pedagógico reforçados', 7),
  ('Utilização efetiva do SIGE para apoio à gestão escolar', 7),
  ('Melhoria da eficiência e eficácia da gestão escolar', 7);

-- Objetivo 7 - Indicadores
INSERT INTO indicadores (indicador, desc_indicador, valor_baseline, unidades_id, objetivo_id)
VALUES 
  ('IND35', 'Nº de Unidades Educativas com órgãos de gestão plenamente funcionais', 71.5, 1, 7),
  ('IND36', 'Nº de dirigentes das Unidades educativas capacitados em elaboração de projeto educativo', NULL, 2, 7),
  ('IND37', 'Nº de Unidades Educativas que adota projeto educativo como instrumento estratégico', 42, 1, 7),
  ('IND38', 'Nº de Unidades Educativas com atividade de supervisão anual', NULL, 1, 7),
  ('IND39', 'Nº de dirigentes e pessoal administrativo capacitados em técnicas do SIGE', NULL, 2, 7),
  ('IND40', '% de escolas com SIGE funcional e dados completos e de qualidade', NULL, 1, 7);

-- Objetivo 7 - Metas 2026
INSERT INTO metas (meta, meta_valor, ano_id, indicador_id)
VALUES 
  ('Meta IND35 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND35')),
  ('Meta IND36 2026', 617, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND36')),
  ('Meta IND37 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND37')),
  ('Meta IND38 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND38')),
  ('Meta IND39 2026', 917, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND39')),
  ('Meta IND40 2026', 100, 2, (SELECT indicador_id FROM indicadores WHERE indicador='IND40'));
