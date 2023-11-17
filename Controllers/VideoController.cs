using System.IO;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Microsoft.AspNetCore.StaticFiles;
using Microsoft.Graph;
using System.Net.Http;
using System.Net.Http.Headers;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Http;
using aspnet_webapp.Services;
using System.Linq;
using Microsoft.AspNetCore.Routing;
using System.ComponentModel.DataAnnotations;

namespace aspnet_webapp.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class VideoController : ControllerBase
    {
        private readonly ILogger _logger;
        private readonly IConfiguration _configuration;
        private readonly IVideoService _videoService;

        public VideoController(IVideoService videoService, IConfiguration configuration,ILogger<VideoController> logger)
        {
            _logger = logger;
            _configuration = configuration;
            _videoService=videoService;
        }

        [HttpGet()]
        public IActionResult GetVideos()
        {
            
            return this.Ok();
        }


        [HttpGet("{name}")]
        public IActionResult GetVideo(string name)
        {
            _logger.LogInformation("Stream video {0}, range {1}", name, this.Request.GetTypedHeaders().Range?.ToString());
            var video=_videoService.GetVideos().FirstOrDefault(v=>v.Id==name);

            if(video==null) {
                return this.NotFound();
            }
            
            return this.PhysicalFile(video.FileName, "video/mp4", true);
        }

    }
}