(function () {

	/**
	 * Creates the angular application.
	 */
	angular.module('app', [
			'ngRoute'
		]);

	/**
	 * Defines the route configuration for the
	 * main application.
	 */
	function Config ($routeProvider, $httpProvider, $locationProvider) {
		$routeProvider
		.when('/', {
			templateUrl: '/static/lgtm.html',
			controller: 'RepoCtrl'
		})
        .when('/:org', {
            templateUrl: '/static/lgtm.html',
            controller: 'RepoCtrl'
        })
        ;

		// Enables html5 mode
		$locationProvider.html5Mode(true);

        // Enables XSRF protection
        $httpProvider.defaults.headers.common['X-CSRF-TOKEN'] = window.STATE_FROM_SERVER.csrf;
	}

    function Noop($rootScope) {}

	angular
		.module('app')
		.config(Config)
        .run(Noop);
})();

(function () {

	function parseRepo() {
	    return function(conf_url) {
			var parts = conf_url.split("/")
			return parts[3]+"/"+parts[4];
	    }
	}

	angular
		.module('app')
		.filter('parseRepo', parseRepo);

})();

(function () {
	function UserService($http) {
        var user_ = window.STATE_FROM_SERVER.user;

		this.current = function() {
			return user_;
		};
	}

	angular
		.module('app')
		.service('user', UserService);
})();

(function () {
	function TeamService($http) {
        var teams_ = window.STATE_FROM_SERVER.teams || [];
        teams_.unshift(window.STATE_FROM_SERVER.user);

		this.list = function() {
            return teams_;
		};

        this.get = function(name) {
            for (var i=0; i<teams_.length; i++) {
                if (teams_[i].login === name) {
                    return teams_[i];
                }
            }
        }
	}

	angular
		.module('app')
		.service('teams', TeamService);
})();

(function () {
	function RepoService($http) {

		this.list = function() {
			return $http.get('/api/user/repos');
		};

		this.post = function(repo, body) {
			return $http.post('/api/repos/'+repo.owner+'/'+repo.name, body);
		};

        this.delete = function(repo) {
			return $http.delete('/api/repos/'+repo.owner+'/'+repo.name);
		};
	}

	angular
		.module('app')
		.service('repos', RepoService);
})();

(function () {
	function RepoCtrl($scope, $routeParams, repos, teams, user) {

        $scope.org = teams.get($routeParams.org || user.current().login);
        $scope.orgs = teams.list();
        $scope.user = user.current();

		repos.list().then(function(payload){
			$scope.repos = payload.data;
			delete $scope.error;
		}).catch(function(err){
			$scope.error = err;
		});

		$scope.activate = function(repo) {
			var index = $scope.repos.indexOf(repo);
			repos.post(repo, {}).then(function(payload){
                delete $scope.repo;
				delete $scope.error;
                $scope.repos[index] = payload.data;
                $scope.saving = false;
			}).catch(function(err){
                delete $scope.repo;
				$scope.error = err;
                $scope.saving = false;
			});
            $scope.saving = true;
		}

        $scope.delete = function(repo) {
            delete repo.id;
			repos.delete(repo).catch(function(err){
				$scope.error = err;
			});
		}

        $scope.changeOrg = function(value) {
            $scope.org = teams.get(value);
        }
		$scope.edit = function(repo) {
            $scope.repo = repo;
		};
        $scope.close = function() {
            delete $scope.repo;
		};
        $scope.saving = false;
	}

	angular
		.module('app')
		.controller('RepoCtrl', RepoCtrl);
})();
