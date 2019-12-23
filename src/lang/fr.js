export default {
  table: {
    date: 'Date',
    create: 'Creation',
    modification: 'Modification',
    role: 'Role',
    add: 'Ajouter',
    status: 'Status',
    title: 'Titre',
    active: 'Actif',
    select: 'Selectionner',
    action: 'Action',
    delete: 'Effacer',
    cancel: 'Annuler',
    save: 'Sauvegarder',
    confirm: 'Sauvegarder',
    update: 'Mettre a jour',
    edit: 'Modifier'
  },
  post: {
    display_time: 'Date',
    display_time_placeholder: 'Sélectionnez',
    title: 'Titre',
    userinfo: 'Auteur',
    author_select: 'Sélectionnez un auteur',
    author_rating: 'Évaluation',
    description: 'Description',
    description_placeholder: 'Entrez une description ici',
    content: 'Content',
    content_placeholder: 'Enter description here',
    cover_image: 'Image'
  },
  plugins: {
    'serverless-cms': {
      button: {
        launch_create: 'Lancer la création',
        launch_create_cluster: 'Lancer la création',
        clear: 'Effacer le formulaire',
        save_state: 'Enregistrer',
        refresh_state: 'Rafraîchir'
      },
      error: {
        not_started: 'La création d\'image pour {title} n\'a pas été démarré.',
        not_started_clusetr: 'La création du cluster {title} n\'a pas été démarré.'
      },
      clustertitle: 'Création du cluster',
      unique_stack_name: 'Le nom doit être unique.',
      stack_state: 'Le status de {name} est {status}.',
      environment: 'Valeur du tag d\'environment ainsi que le préfixe pour les ressources créé',
      'createserverless-cmstitle': 'Création du cluster ECS',
      ecs: 'Informations',
      createecsnodegroup: 'Création de serveurs',
      createecsdiscoverygroup: 'Création du service de discovery',
      task_placement_tags: 'Si défini, les serveurs vont se limiter qu\'aux tâches ayant cette valeur'
    },
    test: {
      amiCreateTitle: 'Création d\'une image personalisé',
      personal_details: 'Détails personnel',
      region_details: 'AMI & Region',
      ami: 'Image ( AMI )'
    }
  },
  error: {
    no_vpc_in_region: 'Aucun VPC\'s disponible dans la region {region}',
    no_subnet_in_region_vpc: 'Aucun résau disponible en {region} pour le vpc {vpc}',
    input: {
      title: 'Vous devez spécifier un nom qui sera public',
      email: 'Vous devez entrer une address courrirel valide'
    }
  },
  user: {
    title: 'Nom Public',
    email: 'Courriel'
  }
}
